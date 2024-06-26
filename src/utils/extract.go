package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Unzip(src string, dest string) error {
	fmt.Println("Extracting Gzip from " + src + " to " + dest)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		if err := os.MkdirAll(dest, 0755); err != nil {
			panic(err)
		}
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, strings.Join(strings.Split(f.Name, "/")[1:], "/"))

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File[1:] {
		if strings.Contains(f.Name, "..") {
			continue
		}
		err := extractAndWriteFile(f)
		if err != nil {
			return err
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

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		if err := os.MkdirAll(dest, 0755); err != nil {
			return err
		}
	}

	gzr, err := gzip.NewReader(r)
	if err != nil {
		r.Close()
		// Check if make command is available
		extract_cmd := exec.Command(
			"tar",
			"-xvf",
			src,
			"-C",
			dest,
		)

		fmt.Println(src, dest)
		out, err := extract_cmd.StdoutPipe()
		if err != nil {
			return err
		}

		if err = extract_cmd.Start(); err != nil {
			return fmt.Errorf("unable to extract the file")
		}
		for {
			tmp := make([]byte, 1024)
			_, err := out.Read(tmp)
			fmt.Print(string(tmp))
			if err != nil {
				break
			}
		}
		return nil
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil || strings.Contains(header.Name, ".."):
			continue
		}

		insideZipPath := ""

		if directExtract {
			insideZipPath = strings.Join(strings.Split(header.Name, "/")[0:], "/")
		} else {
			insideZipPath = strings.Join(strings.Split(header.Name, "/")[1:], "/")
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dest, insideZipPath)

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()

		// if it's a symlink create it
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}

		// if it's a hardlink create it
		case tar.TypeLink:
			if err := os.Link(header.Linkname, target); err != nil {
				return err
			}
		}
	}
}
