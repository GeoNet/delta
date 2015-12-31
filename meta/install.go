package meta

type Install struct {
	Equipment
	Span
}

func (i Install) Less(in Install) bool {
	switch {
	case i.Equipment.Less(in.Equipment):
		return true
	case in.Equipment.Less(i.Equipment):
		return false
	default:
		return i.Span.Less(in.Span)
	}
}
