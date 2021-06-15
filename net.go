package main

import (
	"bytes"
	"io"
	"net/http"

	"golang.org/x/xerrors"
)

func Get(url string, progress bool) (*bytes.Buffer, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf("request error[%s]: %w", url, err)
	}
	defer resp.Body.Close()

	var b bytes.Buffer

	if progress {
		p := NewProgressWriter(&b, resp.ContentLength)
		p.Event = PrefixProgressFunc("Download")
		_, err = p.Copy(resp.Body)
		if err != nil {
			return nil, xerrors.Errorf("io.Copy(): %w", err)
		}
	} else {
		_, err = io.Copy(&b, resp.Body)
		if err != nil {
			return nil, xerrors.Errorf("io.Copy(): %w", err)
		}
	}

	return &b, nil
}
