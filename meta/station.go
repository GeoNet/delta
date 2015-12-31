package meta

type Station struct {
	Reference
	Point
	Span

	Notes string
}

type Stations []Station

func (s Stations) Len() int           { return len(s) }
func (s Stations) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Stations) Less(i, j int) bool { return s[i].Reference.less(s[j].Reference) }
