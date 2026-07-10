package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestUnGzipExtractsFilesAndSafeRelativeSymlinks(t *testing.T) {
	archivePath := filepath.Join(t.TempDir(), "node.tar.gz")
	writeTarGzip(t, archivePath, []tar.Header{
		{Name: "node-v1/lib/npm-cli.js", Mode: 0644, Size: int64(len("content")), Typeflag: tar.TypeReg},
		{Name: "node-v1/bin/npm", Mode: 0777, Typeflag: tar.TypeSymlink, Linkname: "../lib/npm-cli.js"},
	}, [][]byte{[]byte("content"), nil})
	destination := filepath.Join(t.TempDir(), "node")

	if err := UnGzip(archivePath, destination, false); err != nil {
		t.Fatalf("extract archive: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(destination, "bin", "npm"))
	if err != nil {
		t.Fatalf("read through extracted symlink: %v", err)
	}
	if string(content) != "content" {
		t.Fatalf("unexpected extracted content %q", content)
	}
}

func TestUnGzipRejectsParentDirectoryTraversal(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "node.tar.gz")
	writeTarGzip(t, archivePath, []tar.Header{
		{Name: "node-v1/../../escaped", Mode: 0644, Size: int64(len("malicious")), Typeflag: tar.TypeReg},
	}, [][]byte{[]byte("malicious")})
	destination := filepath.Join(root, "nested", "node")

	if err := UnGzip(archivePath, destination, false); err == nil {
		t.Fatal("expected traversal entry to be rejected")
	}
	if _, err := os.Stat(filepath.Join(root, "escaped")); !os.IsNotExist(err) {
		t.Fatalf("traversal created a file outside destination: %v", err)
	}
}

func TestUnGzipRejectsEscapingSymlink(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "node.tar.gz")
	writeTarGzip(t, archivePath, []tar.Header{
		{Name: "node-v1/bin/escape", Mode: 0777, Typeflag: tar.TypeSymlink, Linkname: "../../../escaped"},
	}, [][]byte{nil})
	destination := filepath.Join(root, "node")

	if err := UnGzip(archivePath, destination, false); err == nil {
		t.Fatal("expected escaping symlink to be rejected")
	}
}

func TestUnzipExtractsFilesAndRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	validArchive := filepath.Join(root, "node.zip")
	writeZip(t, validArchive, map[string]string{"node-v1/bin/node.exe": "binary"})
	destination := filepath.Join(root, "valid")
	if err := Unzip(validArchive, destination); err != nil {
		t.Fatalf("extract valid zip: %v", err)
	}
	if content, err := os.ReadFile(filepath.Join(destination, "bin", "node.exe")); err != nil || string(content) != "binary" {
		t.Fatalf("unexpected extracted zip content %q: %v", content, err)
	}

	maliciousArchive := filepath.Join(root, "malicious.zip")
	writeZip(t, maliciousArchive, map[string]string{"node-v1/../../escaped": "malicious"})
	if err := Unzip(maliciousArchive, filepath.Join(root, "nested", "node")); err == nil {
		t.Fatal("expected zip traversal entry to be rejected")
	}
	if _, err := os.Stat(filepath.Join(root, "escaped")); !os.IsNotExist(err) {
		t.Fatalf("zip traversal created a file outside destination: %v", err)
	}
}

func writeTarGzip(t *testing.T, archivePath string, headers []tar.Header, contents [][]byte) {
	t.Helper()
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("create tarball: %v", err)
	}
	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)
	for index := range headers {
		header := headers[index]
		if err := tarWriter.WriteHeader(&header); err != nil {
			t.Fatalf("write tar header: %v", err)
		}
		if len(contents[index]) > 0 {
			if _, err := tarWriter.Write(contents[index]); err != nil {
				t.Fatalf("write tar content: %v", err)
			}
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatalf("close tar writer: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatalf("close gzip writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close tarball: %v", err)
	}
}

func writeZip(t *testing.T, archivePath string, entries map[string]string) {
	t.Helper()
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	writer := zip.NewWriter(file)
	for name, content := range entries {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatalf("create zip entry: %v", err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			t.Fatalf("write zip entry: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
}
