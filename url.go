package main

import (
	"encoding/xml"
	"fmt"
	"runtime"
	"sort"

	"golang.org/x/xerrors"
)

func getURLs(v *Version) ([]string, error) {

	xmlURL := fmt.Sprintf(XMLURL, v.Major)

	fmt.Println("Version List URL:", xmlURL)

	b, err := Get(xmlURL, false)
	if err != nil {
		return nil, xerrors.Errorf("Get() error: %w", err)
	}

	var cdxml XMLResponse
	err = xml.Unmarshal(b.Bytes(), &cdxml)
	if err != nil {
		return nil, xerrors.Errorf("xml.Unmarshal() error: %w", err)
	}

	ps := cdxml.CommonPrefixes
	if len(ps) <= 0 {
		return nil, fmt.Errorf("Version Not Found")
	}

	arch := getArch()

	if arch == MacM1Arch {
		if !v.NotM1Support() {
			return nil, fmt.Errorf("Not M1 Support jVersion[%v]", v)
		}
	}

	fn := fmt.Sprintf(ZIPFileName, arch)

	prefixes := cdxml.CommonPrefixes
	sort.Slice(prefixes, func(i, j int) bool {

		p1 := prefixes[i].Prefix
		p2 := prefixes[j].Prefix

		v1 := NewVersion(p1[:len(p1)-1])
		v2 := NewVersion(p2[:len(p2)-1])

		if v1.compare(v2) > 0 {
			return true
		}
		return false
	})

	rtn := make([]string, len(ps))
	for idx, pre := range prefixes {
		rtn[idx] = fmt.Sprintf("%s%s%s", RootURL, pre.Prefix, fn)
	}
	return rtn, nil
}

func getArch() string {

	switch runtime.GOOS {
	case "windows":
		return WindowsArch
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return MacM1Arch
		}
		return MacArch
	case "linux":
		return LinuxArch
	}
	return ""
}
