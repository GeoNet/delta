package meta

import (
	"sort"
)

type Network struct {
	Code       string `csv:"Internal Code",`
	Map        string `csv:"External Code",`
	Name       string `csv:"Network Name",`
	Restricted bool   `csv:"Restricted Status",`
}

type Networks []Network

func (s Networks) Len() int           { return len(s) }
func (s Networks) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Networks) Less(i, j int) bool { return s[i].Code < s[j].Code }

func (s Networks) list()      {}
func (s Networks) sort() list { sort.Sort(s); return s }
