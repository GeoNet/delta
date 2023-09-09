package ntrip

import (
	"path/filepath"
)

const (
	AliasesFile = "aliases.csv"
	FormatsFile = "formats.csv"
	ModelsFile  = "models.csv"
	MountsFile  = "mounts.csv"
	UsersFile   = "users.csv"
)

// Caster holds the decoded config data.
type Caster struct {
	aliases []Alias
	formats map[string][]string
	models  map[string]string
	mounts  []Mount
	users   []User
}

// NewCaster returns a Caster pointer after reading all expected config files.
func NewCaster(common, input string) (*Caster, error) {

	formats, err := ReadFormats(filepath.Join(common, FormatsFile))
	if err != nil {
		return nil, err
	}

	models, err := ReadModels(filepath.Join(common, ModelsFile))
	if err != nil {
		return nil, err
	}

	aliases, err := ReadAliases(filepath.Join(input, AliasesFile))
	if err != nil {
		return nil, err
	}

	mounts, err := ReadMounts(filepath.Join(input, MountsFile))
	if err != nil {
		return nil, err
	}

	users, err := ReadUsers(filepath.Join(input, UsersFile))
	if err != nil {
		return nil, err
	}

	caster := Caster{
		aliases: aliases,
		formats: formats,
		models:  models,
		mounts:  mounts,
		users:   users,
	}

	return &caster, nil
}

// Aliases returns all configured aliases as a slice.
func (c *Caster) Aliases() []Alias {
	return append([]Alias{}, c.aliases...)
}

// Format looks up the internal config for the given format details.
func (c *Caster) Format(details string) ([]string, bool) {
	f, ok := c.formats[details]
	return f, ok
}

// Model looks up the internal config for the given model details.
func (c *Caster) Model(model string) (string, bool) {
	m, ok := c.models[model]
	return m, ok
}

// Mounts returns all configured mounts as a slice.
func (c *Caster) Mounts() []Mount {
	return append([]Mount{}, c.mounts...)
}

// Users returns all configured users as a slice.
func (c *Caster) Users() []User {
	return append([]User{}, c.users...)
}
