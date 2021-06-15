package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/xerrors"
)

var download bool
var version string
var dir string

func init() {
	flag.BoolVar(&download, "d", false, "driver download")
	flag.StringVar(&version, "v", "", "specify version")
}

func main() {

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		if !download {
			fmt.Println("the path is specified even though i do not download it.")
		}
		dir = args[0]
	}

	err := run()
	if err != nil {
		fmt.Printf("run error:\n%+v\n", err)
		os.Exit(1)
	}
	fmt.Println("Success")
}

func run() error {

	if dir != "" {
		if _, err := os.Stat(dir); err != nil {
			return xerrors.Errorf("os.Stat() error: %w", err)
		}
	}

	var err error
	var ver *Version
	if version != "" {
		ver = NewVersion(version)
	} else {
		ver, err = GetChromeVersion()
		if err != nil {
			return xerrors.Errorf("GetChromeVersion() error: %w", err)
		}
		fmt.Println("Now Install Chrome Version", ver)
	}

	if ver.IsZero() {
		return fmt.Errorf("version parse error[%s]", ver.Src)
	}

	if ver.NotSupport() {
		return fmt.Errorf("Not Support Version[%s]", ver.Src)
	}

	urls, err := getURLs(ver)
	if err != nil {
		return xerrors.Errorf("printDownloadURL() error: %w", err)
	}

	if len(urls) <= 0 {
		return xerrors.Errorf("Not Found Version")
	}

	if download {
		url := urls[0]
		err = DownloadAndWrite(dir, url)
		if err != nil {
			return xerrors.Errorf("DownloadAndWrite() error: %w", err)
		}
	} else {
		fmt.Println("Download URL List")
		for idx, url := range urls {
			fmt.Printf("URL%d:[%s]\n", idx+1, url)
		}
	}

	return nil
}
