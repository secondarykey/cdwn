package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/xerrors"
)

func DownloadAndWrite(dir string, url string) error {

	fmt.Println("ZIP File Download", url)
	buf, err := Get(url, true)
	if err != nil {
		return xerrors.Errorf("stream.Get() error: %w", err)
	}
	fmt.Println()

	name := FileName
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	path := name
	if dir != "" {
		path = filepath.Join(dir, name)
	}

	fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 755)
	if err != nil {
		return xerrors.Errorf("os.Create() error: %w", err)
	}
	defer fp.Close()

	b := bytes.NewReader(buf.Bytes())
	err = CopyZIP(fp, b)
	if err != nil {
		return xerrors.Errorf("CopyZIP() error: %w", err)
	}

	fmt.Println("Create Chrome Driver", path)
	return nil
}

func CopyZIP(w io.Writer, r *bytes.Reader) error {

	z, err := zip.NewReader(r, int64(r.Len()))
	if err != nil {
		return xerrors.Errorf("zip.NewReader() error: %w", err)
	}

	if len(z.File) > 1 {
		return fmt.Errorf("multiple file not support")
	}

	f := z.File[0]
	fp, err := f.Open()
	if err != nil {
		return xerrors.Errorf("zip file Open() error: %w", err)
	}
	defer fp.Close()

	info := f.FileHeader.FileInfo()
	prog := NewProgressWriter(w, info.Size())

	fmt.Println("Uncompressed ZIP File")
	_, err = prog.Copy(fp)
	if err != nil {
		return xerrors.Errorf("io.Open() error: %w", err)
	}
	fmt.Println()
	return nil
}
