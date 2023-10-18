package main

import (
	"io"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v2"
)

// MountConfig represents a set if ntripcaster mount point configuration details.
type MountConfig struct {
	Mount     string `yaml:"mount"`
	Mark      string `yaml:"mark"`
	Name      string `yaml:"name"`
	Latitude  string `yaml:"latitude"`
	Longitude string `yaml:"longitude"`
	Country   string `yaml:"country"`
	Format    string `yaml:"format"`
}

func (m *MountConfig) Less(mount MountConfig) bool {
	return m.Mount < mount.Mount
}

// Config represents a set if bnc mount details.
type Config struct {
	Mounts []MountConfig `yaml:"gns_bnc::mounts"`
}

func (c *Config) Sort() {

	sort.Slice(c.Mounts, func(i, j int) bool {
		return c.Mounts[i].Less(c.Mounts[j])
	})

}

func (c *Config) Marshal() ([]byte, error) {

	data, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Config) Write(wr io.Writer) error {

	data, err := c.Marshal()
	if err != nil {
		return err
	}

	if _, err := wr.Write(data); err != nil {
		return err
	}

	return nil
}

func (c *Config) WriteFile(path string) error {

	data, err := c.Marshal()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}

	return nil
}
