package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Unzip(src string, dest string) error {
	fmt.Println("Extracting Zip from " + src + " to " + dest)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := ensureDestinationDirectory(dest); err != nil {
		return err
	}

	for _, entry := range r.File {
		relative, err := archiveRelativePath(entry.Name, false)
		if err != nil {
			return err
		}
		if relative == "" {
			continue
		}
		if entry.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("zip entry %q is a symbolic link", entry.Name)
		}

		target, err := safeArchiveTarget(dest, relative)
		if err != nil {
			return err
		}
		if entry.FileInfo().IsDir() {
			if err := ensureSafeParents(dest, target); err != nil {
				return err
			}
			if err := os.Mkdir(target, entry.Mode().Perm()); err != nil && !os.IsExist(err) {
				return err
			}
			continue
		}

		if err := ensureSafeParents(dest, filepath.Dir(target)); err != nil {
			return err
		}
		if err := rejectSymlink(target); err != nil {
			return err
		}

		rc, err := entry.Open()
		if err != nil {
			return err
		}
		output, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, entry.Mode().Perm())
		if err != nil {
			rc.Close()
			return err
		}
		_, copyErr := io.Copy(output, rc)
		closeOutputErr := output.Close()
		closeInputErr := rc.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeOutputErr != nil {
			return closeOutputErr
		}
		if closeInputErr != nil {
			return closeInputErr
		}
	}

	return nil
}

func UnGzip(src string, dest string, directExtract bool) error {
	fmt.Println("Extracting Gzip from " + src + " to " + dest)
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("open gzip stream: %w", err)
	}
	defer gzr.Close()

	if err := ensureDestinationDirectory(dest); err != nil {
		return err
	}

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		relative, err := archiveRelativePath(header.Name, directExtract)
		if err != nil {
			return err
		}
		if relative == "" {
			continue
		}
		target, err := safeArchiveTarget(dest, relative)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := ensureSafeParents(dest, target); err != nil {
				return err
			}
			if err := os.Mkdir(target, os.FileMode(header.Mode).Perm()); err != nil && !os.IsExist(err) {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := ensureSafeParents(dest, filepath.Dir(target)); err != nil {
				return err
			}
			if err := rejectSymlink(target); err != nil {
				return err
			}
			output, err := os.OpenFile(
				target,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				os.FileMode(header.Mode).Perm(),
			)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(output, tr)
			closeErr := output.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		case tar.TypeSymlink:
			if err := ensureSafeParents(dest, filepath.Dir(target)); err != nil {
				return err
			}
			linkTarget := filepath.FromSlash(header.Linkname)
			if filepath.IsAbs(linkTarget) {
				return fmt.Errorf("archive symlink %q has an absolute target", header.Name)
			}
			if _, err := safeArchiveTarget(dest, filepath.Join(filepath.Dir(relative), linkTarget)); err != nil {
				return fmt.Errorf("unsafe archive symlink %q: %w", header.Name, err)
			}
			if err := rejectExistingPath(target); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		case tar.TypeLink:
			linkRelative, err := archiveRelativePath(header.Linkname, directExtract)
			if err != nil {
				return fmt.Errorf("unsafe archive hardlink %q: %w", header.Name, err)
			}
			linkTarget, err := safeArchiveTarget(dest, linkRelative)
			if err != nil {
				return err
			}
			if err := ensureSafeParents(dest, filepath.Dir(target)); err != nil {
				return err
			}
			if err := rejectExistingPath(target); err != nil {
				return err
			}
			if err := os.Link(linkTarget, target); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported tar entry type %d for %q", header.Typeflag, header.Name)
		}
	}
}

func archiveRelativePath(name string, directExtract bool) (string, error) {
	name = strings.ReplaceAll(name, "\\", "/")
	cleaned := path.Clean(name)
	if cleaned == "." {
		return "", nil
	}
	if path.IsAbs(cleaned) {
		return "", fmt.Errorf("archive entry %q is absolute", name)
	}

	parts := strings.Split(cleaned, "/")
	for _, part := range parts {
		if part == ".." {
			return "", fmt.Errorf("archive entry %q escapes its destination", name)
		}
	}
	if !directExtract {
		if len(parts) == 1 {
			return "", nil
		}
		parts = parts[1:]
	}

	relative := filepath.FromSlash(strings.Join(parts, "/"))
	if filepath.IsAbs(relative) || filepath.VolumeName(relative) != "" {
		return "", fmt.Errorf("archive entry %q is absolute", name)
	}
	return relative, nil
}

func safeArchiveTarget(dest, relative string) (string, error) {
	destination := filepath.Clean(dest)
	target := filepath.Clean(filepath.Join(destination, relative))
	inside, err := filepath.Rel(destination, target)
	if err != nil || inside == ".." || strings.HasPrefix(inside, ".."+string(os.PathSeparator)) || filepath.IsAbs(inside) {
		return "", fmt.Errorf("archive path %q escapes destination %q", relative, dest)
	}
	return target, nil
}

func ensureDestinationDirectory(dest string) error {
	metadata, err := os.Lstat(dest)
	if err == nil {
		if metadata.Mode()&os.ModeSymlink != 0 || !metadata.IsDir() {
			return fmt.Errorf("archive destination %q is not a directory", dest)
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(dest, 0755)
}

func ensureSafeParents(dest, target string) error {
	destination := filepath.Clean(dest)
	relative, err := filepath.Rel(destination, filepath.Clean(target))
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("path %q escapes archive destination %q", target, dest)
	}

	current := destination
	for _, component := range strings.Split(relative, string(os.PathSeparator)) {
		if component == "" || component == "." {
			continue
		}
		current = filepath.Join(current, component)
		metadata, err := os.Lstat(current)
		if os.IsNotExist(err) {
			if err := os.Mkdir(current, 0755); err != nil && !os.IsExist(err) {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if metadata.Mode()&os.ModeSymlink != 0 || !metadata.IsDir() {
			return fmt.Errorf("archive parent %q is not a safe directory", current)
		}
	}
	return nil
}

func rejectSymlink(target string) error {
	metadata, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if metadata.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("archive target %q is a symbolic link", target)
	}
	return nil
}

func rejectExistingPath(target string) error {
	_, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("archive target %q already exists", target)
}
