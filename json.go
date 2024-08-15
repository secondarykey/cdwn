package main

import (
	"encoding/json"
	"fmt"
	"runtime"

	"golang.org/x/xerrors"
)

type ChromeJSON struct {
	Version   string    `json:"version"`
	Revision  string    `json:"revision"`
	Downloads Downloads `json:"downloads"`
}

type Downloads struct {
	Drivers []DownloadURL `json:"chromedriver"`
	//chrome
	//chrome-headless-shell
}

type DownloadURL struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

func getURL(v string) (string, error) {

	const url = "https://googlechromelabs.github.io/chrome-for-testing/%s.json"
	fn := fmt.Sprintf(url, v)

	fmt.Println("URL:", fn)

	json, err := getDownloadJson(fmt.Sprintf(url, v))
	if err != nil {
		return "", xerrors.Errorf("getDownloadJson() error: %w", err)
	}

	drivers := json.Downloads.Drivers
	if len(drivers) <= 0 {
		return "", xerrors.Errorf("Version not found.[%s]", v)
	}

	arch := getArch()
	if arch == "" {
		return "", xerrors.Errorf("Not Support Architect[%s:%s]", runtime.GOOS, runtime.GOARCH)
	}

	fmt.Println("Architect:", arch)
	for _, d := range drivers {
		if d.Platform == arch {
			return d.URL, nil
		}
	}
	return "", fmt.Errorf("Not Found Platform[%d]", arch)
}

func getDownloadJson(url string) (*ChromeJSON, error) {

	resp, err := Get(url, false)
	if err != nil {
		return nil, xerrors.Errorf("Get() error: %w", err)
	}

	var obj ChromeJSON
	err = json.Unmarshal(resp.Bytes(), &obj)
	if err != nil {
		return nil, xerrors.Errorf("json.Unmarshal() error: %w", err)
	}

	return &obj, nil
}

func getArch() string {

	arch := runtime.GOARCH
	switch runtime.GOOS {
	case "windows":
		if arch == "arm64" || arch == "amd64" {
			return Windows64Arch
		}
		return Windows32Arch
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return MacM1Arch
		}
		return MacArch
	case "linux":
		return LinuxArch
	default:
	}
	return ""
}
