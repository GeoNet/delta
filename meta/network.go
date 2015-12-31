package meta

type Network struct {
	NetworkCode  string
	ExternalCode string
	Description  string
	Restricted   bool
}

type Networks []Network

func (n Networks) Len() int           { return len(n) }
func (n Networks) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Networks) Less(i, j int) bool { return n[i].NetworkCode < n[j].NetworkCode }
