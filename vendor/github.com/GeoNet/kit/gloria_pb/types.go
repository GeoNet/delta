package gloria_pb

// Additional types for working with Gloria protobufs.  This file is not automatically generated.

type MarksPriority []Mark

func (m MarksPriority) Len() int {
	return len(m)
}

func (m MarksPriority) Less(i, j int) bool {
	return m[i].GetDownload().GetPriority() < m[j].GetDownload().GetPriority()
}

func (m MarksPriority) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
