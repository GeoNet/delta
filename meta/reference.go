package meta

type Reference struct {
	Code    string
	Network string
	Name    string
}

func (r Reference) Less(ref Reference) bool { return r.Code < ref.Code }
