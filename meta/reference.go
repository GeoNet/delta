package meta

type Reference struct {
	Code    string
	Network string
	Name    string
}

func (r Reference) less(ref Reference) bool { return r.Code < ref.Code }
