package main

const (
	Query   = "HKEY_CURRENT_USER\\Software\\Google\\Chrome\\BLBeacon"
	RootURL = "https://chromedriver.storage.googleapis.com/"
	XMLURL  = RootURL + "?delimiter=/&prefix=%d"
)

const (
	ZIPFileName = "chromedriver_%s.zip"
	WindowsArch = "win32"    //support from v2.0
	MacArch     = "mac64"    //70.0.3538.16
	MacM1Arch   = "mac64_m1" //87.0.4280.88 support
	LinuxArch   = "linux64"  //70.0.3538.16 support

	FileName = "chromedriver"
)
