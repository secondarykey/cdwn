package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/xerrors"
)

func DownloadAndWrite(dir string, url string) error {

	fmt.Println("ZIP File", url)

	buf, err := Get(url, true)
	if err != nil {
		return xerrors.Errorf("Get() error: %w", err)
	}
	fmt.Println()

	b := bytes.NewReader(buf.Bytes())
	err = CopyZIP(dir, b)
	if err != nil {
		return xerrors.Errorf("CopyZIP() error: %w", err)
	}

	return nil
}

func Extract[From, To any](s []From, f func(From) To) []To {
	res := make([]To, len(s))
	for i, v := range s {
		res[i] = f(v)
	}
	return res
}

func CopyZIP(out string, r *bytes.Reader) error {

	name := FileName
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	z, err := zip.NewReader(r, int64(r.Len()))
	if err != nil {
		return xerrors.Errorf("zip.NewReader() error: %w", err)
	}

	f := findTargetZipFile(z.File, name)
	if f == nil {
		n := Extract(z.File, func(f *zip.File) string { return f.FileHeader.Name })
		return fmt.Errorf("[%s] file Not Found[%s]", strings.Join(n, ","))
	}

	zp, err := f.Open()
	if err != nil {
		return xerrors.Errorf("zip file Open() error: %w", err)
	}
	defer zp.Close()
	info := f.FileHeader.FileInfo()

	p := filepath.Join(out, name)
	fp, err := os.Create(p)
	if err != nil {
		return xerrors.Errorf("os.Create() error: %w", err)
	}
	defer fp.Close()

	prog := NewProgressWriter(fp, info.Size())
	prog.Event = PrefixProgressFunc("Uncompress")
	_, err = prog.Copy(zp)
	if err != nil {
		return xerrors.Errorf("io.Open() error: %w", err)
	}
	fmt.Println()
	fmt.Printf("Create Chrome Driver[%s]\n", p)

	return nil
}

func findTargetZipFile(files []*zip.File, name string) *zip.File {
	for _, f := range files {
		if sameZipFileName(f, name) {
			return f
		}
	}
	return nil
}

func sameZipFileName(f *zip.File, name string) bool {
	n := f.FileHeader.Name
	idx := strings.LastIndex(n, name)
	if idx == -1 {
		return false
	}

	sz := len(n)
	l := len(name)

	if idx == (sz - l) {
		return true
	}
	return false
}
