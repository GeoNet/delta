package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func readLastSiteLog(dir, code string) (*SiteLog, error) {
	var files []string

	// check existing files
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".xml" {
			return nil
		}

		if !strings.Contains(filepath.Base(path), "_") {
			return nil
		}

		if !strings.HasPrefix(filepath.Base(path), strings.ToLower(code)) {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return filepath.Base(files[i]) < filepath.Base(files[j])
	})

	if n := len(files); n > 0 {
		raw, err := ioutil.ReadFile(files[n-1])
		if err != nil {
			return nil, err
		}

		var last SiteLogInput
		if err := xml.Unmarshal(raw, &last); err != nil {
			return nil, err
		}

		return last.SiteLog(), nil
	}

	return nil, nil
}
