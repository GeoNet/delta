package main

import (
	"io"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v2"
)

// Store ntripcaster configuration as requred for a heira yaml file

// ClientMountConfig represents a set if ntripcaster user configuration details.
type UserConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (u *UserConfig) Less(user UserConfig) bool {
	return u.Username < user.Username
}

// ClientMountConfig represents a set if ntripcaster group configuration details.
type GroupConfig struct {
	Group string   `yaml:"group"`
	Users []string `yaml:"users"`
}

func (g *GroupConfig) Less(group GroupConfig) bool {
	return g.Group < group.Group
}

func (g *GroupConfig) Sort() {
	sort.Strings(g.Users)
}

// ClientMountConfig represents a set if ntripcaster client mount configuration details.
type ClientMountConfig struct {
	Mount  string   `yaml:"mount"`
	Groups []string `yaml:"groups"`
}

func (m *ClientMountConfig) Less(mount ClientMountConfig) bool {
	return m.Mount < mount.Mount
}

func (m *ClientMountConfig) Sort() {
	sort.Strings(m.Groups)
}

// MountConfig represents a set if ntripcaster mount point configuration details.
type MountConfig struct {
	Mount      string `yaml:"mount"`
	Mark       string `yaml:"mark"`
	Name       string `yaml:"name"`
	Latitude   string `yaml:"latitude"`
	Longitude  string `yaml:"longitude"`
	Country    string `yaml:"country"`
	Format     string `yaml:"format"`
	Details    string `yaml:"details"`
	Navigation string `yaml:"navigation"`
	Model      string `yaml:"model"`
	User       string `yaml:"user,omitempty"`
	Address    string `yaml:"address"`
}

func (m *MountConfig) Less(mount MountConfig) bool {
	return m.Mount < mount.Mount
}

type AliasConfig struct {
	Alias string `yaml:"alias"`
	Mount string `yaml:"mount"`
}

func (a *AliasConfig) Less(alias AliasConfig) bool {
	return a.Alias < alias.Alias
}

// Config represents a set if ntripcaster configuration details.
type Config struct {
	Users        []UserConfig        `yaml:"gns_ntripcaster::users"`
	Groups       []GroupConfig       `yaml:"gns_ntripcaster::groups"`
	ClientMounts []ClientMountConfig `yaml:"gns_ntripcaster::clientmounts"`
	Mounts       []MountConfig       `yaml:"gns_ntripcaster::mounts"`
	Aliases      []AliasConfig       `yaml:"gns_ntripcaster::aliases"`
}

func (c *Config) Sort() {

	sort.Slice(c.Users, func(i, j int) bool {
		return c.Users[i].Less(c.Users[j])
	})

	for _, g := range c.Groups {
		g.Sort()
	}

	sort.Slice(c.Groups, func(i, j int) bool {
		return c.Groups[i].Less(c.Groups[j])
	})

	for _, m := range c.ClientMounts {
		m.Sort()
	}

	sort.Slice(c.ClientMounts, func(i, j int) bool {
		return c.ClientMounts[i].Less(c.ClientMounts[j])
	})

	sort.Slice(c.Mounts, func(i, j int) bool {
		return c.Mounts[i].Less(c.Mounts[j])
	})

	sort.Slice(c.Aliases, func(i, j int) bool {
		return c.Aliases[i].Less(c.Aliases[j])
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
