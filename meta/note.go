package meta

import (
	"sort"
	"strings"
)

const (
	noteCode = iota
	noteNetwork
	noteEntry
	noteLast
)

var noteHeaders Header = map[string]int{
	"Code":    noteCode,
	"Network": noteNetwork,
	"Entry":   noteEntry,
}

var NoteTable Table = Table{
	name:    "Note",
	headers: noteHeaders,
	primary: []string{"Code", "Network"},
	foreign: map[string][]string{
		"Network": {"Network"},
	},
}

type Note struct {
	Code    string
	Network string
	Entry   string
}

type NoteList []Note

func (n NoteList) Len() int      { return len(n) }
func (n NoteList) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n NoteList) Less(i, j int) bool {
	switch {
	case n[i].Code < n[j].Code:
		return true
	case n[i].Code > n[j].Code:
		return false
	case n[i].Network < n[j].Network:
		return true
	case n[i].Network > n[j].Network:
		return false
	case n[i].Entry < n[j].Entry:
		return true
	default:
		return false
	}
}

func (n NoteList) encode() [][]string {
	var data [][]string

	data = append(data, noteHeaders.Columns())

	for _, l := range n {
		data = append(data, []string{
			strings.TrimSpace(l.Code),
			strings.TrimSpace(l.Network),
			strings.TrimSpace(l.Entry),
		})
	}
	return data
}

func (n *NoteList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var notes []Note

	fields := noteHeaders.Fields(data[0])
	for _, v := range data[1:] {
		d := fields.Remap(v)

		notes = append(notes, Note{
			Code:    strings.TrimSpace(d[noteCode]),
			Network: strings.TrimSpace(d[noteNetwork]),
			Entry:   strings.TrimSpace(d[noteEntry]),
		})
	}

	*n = NoteList(notes)

	return nil
}

func LoadNotes(path string) ([]Note, error) {
	var v []Note

	if err := LoadList(path, (*NoteList)(&v)); err != nil {
		return nil, err
	}

	sort.Sort(NoteList(v))

	return v, nil
}
