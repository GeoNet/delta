package main

import (
	"path/filepath"

	"github.com/GeoNet/delta/meta"
)

func NetworkMap(network string) (map[string]meta.Network, error) {

	networkMap := make(map[string]meta.Network)

	var n meta.NetworkList
	if err := meta.LoadList(filepath.Join(network, "networks.csv"), &n); err != nil {
		return nil, err
	}

	for _, v := range n {
		networkMap[v.NetworkCode] = v
	}

	return networkMap, nil
}
