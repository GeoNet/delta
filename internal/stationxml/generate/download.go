package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
)

// download the remote url contents, an insecure option is available for problem sites.
func download(url string, insecure bool) ([]byte, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, //nolint:gosec // needed until fdsn.org passed checks
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
