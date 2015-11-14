package meta

import (
	"sort"
)

type Network struct {
	Code        string `csv:"Network Code"`
	External    string `csv:"External Code"`
	Description string `csv:"Description"`
	Restricted  bool   `csv:"Restricted"`
}

type Networks []Network

func (n Networks) Len() int           { return len(n) }
func (n Networks) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Networks) Less(i, j int) bool { return n[i].Code < n[j].Code }

func (n Networks) List()      {}
func (n Networks) Sort() List { sort.Sort(n); return n }
