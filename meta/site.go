package meta

type Site struct {
	Point

	StationCode  string
	LocationCode string
}

func (s Site) less(site Site) bool {
	switch {
	case s.StationCode < site.StationCode:
		return true
	case s.StationCode > site.StationCode:
		return false
	case s.LocationCode < site.LocationCode:
		return true
	case s.LocationCode > site.LocationCode:
		return false
	default:
		return false
	}
}

type Sites []Site

func (s Sites) Len() int           { return len(s) }
func (s Sites) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Sites) Less(i, j int) bool { return s[i].less(s[j]) }
