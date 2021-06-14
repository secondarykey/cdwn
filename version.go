package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

func GetChromeVersion() (*Version, error) {

	var err error
	line := ""
	if runtime.GOOS == "windows" {
		line, err = getWindowsChromeVersion()
	} else {
		line, err = getChromeVersion()
	}
	if err != nil {
		return nil, xerrors.Errorf("getChromeVersion error: %w", err)
	}

	v, err := parseVersion(line)
	if err != nil {
		return nil, xerrors.Errorf("parseVersion() error: %w", err)
	}
	return v, nil
}

func getWindowsChromeVersion() (string, error) {
	out, err := exec.Command("reg", "query", Query, "/v", "version").Output()
	if err != nil {
		return "", xerrors.Errorf("register command error: %w", err)
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		buf := strings.Trim(line, " \r\t")
		if buf == "" {
			continue
		}
		if strings.Index(buf, "version") == -1 {
			continue
		}
		return buf, nil
	}
	return "", fmt.Errorf("Not Found Version line")
}

func getChromeVersion() (string, error) {
	out, err := exec.Command("google-chrome", "--version").Output()
	if err != nil {
		return "", xerrors.Errorf("Chrome Version error: %w", err)
	}
	return string(out), nil
}

func parseVersion(line string) (*Version, error) {
	versions := strings.Fields(line)
	if len(versions) < 3 {
		return nil, fmt.Errorf("Required 3 column. -> [%s]\n", line)
	}
	num := versions[2]
	return NewVersion(num), nil
}

func parseNum(buf string) int {
	if v, err := strconv.Atoi(buf); err == nil {
		return v
	}
	return -1
}

//https://www.chromium.org/developers/version-numbers
type Version struct {
	Major int
	Minor int
	Build int
	Patch int
	Src   string
}

func NewVersion(buf string) *Version {
	var v Version
	v.Src = buf

	v.Major, v.Minor, v.Build, v.Patch = -1, -1, -1, -1
	vs := strings.Split(buf, ".")
	switch l := len(vs); {
	case l >= 4:
		v.Patch = parseNum(vs[3])
		fallthrough
	case l == 3:
		v.Build = parseNum(vs[2])
		fallthrough
	case l == 2:
		v.Minor = parseNum(vs[1])
		fallthrough
	case l == 1:
		v.Major = parseNum(vs[0])
	}
	return &v
}

func (v1 *Version) compare(v2 *Version) int {
	if v1.Major < v2.Major {
		return -1
	} else if v1.Major > v2.Major {
		return 1
	} else if v1.Minor < v2.Minor {
		return -1
	} else if v1.Minor > v2.Minor {
		return 1
	} else if v1.Build < v2.Build {
		return -1
	} else if v1.Build > v2.Build {
		return 1
	} else if v1.Patch < v2.Patch {
		return -1
	} else if v1.Patch > v2.Patch {
		return 1
	}
	return 0
}

func (v *Version) IsZero() bool {
	if v.Major == -1 && v.Minor == -1 &&
		v.Build == -1 && v.Patch == -1 {
		return true
	}
	return false
}

func (v *Version) String() string {
	return v.Src
}
