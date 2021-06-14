package main

const (
	Query   = "HKEY_CURRENT_USER\\Software\\Google\\Chrome\\BLBeacon"
	RootURL = "https://chromedriver.storage.googleapis.com/"
	XMLURL  = RootURL + "?delimiter=/&prefix=%d"
)

const (
	ZIPFileName = "chromedriver_%s.zip"
	WindowsArch = "win32"    //support from v2.0
	Mac32Arch   = "mac32"    //support up to v2.9
	Mac64Arch   = "mac64"    //70.0.3538.16
	Mac64M1Arch = "mac64_m1" //87.0.4280.88 support
	Linux64Arch = "linux64"  //70.0.3538.16 support
	Linux32Arch = "linux64"  //support up to version 2.9

	FileName = "chromedriver"
)
