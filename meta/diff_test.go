package meta_test

import (
	"fmt"
	"strings"
)

// simple debugging helper function
func diff(s1, s2 []byte) string {

	l1 := strings.Split(
		strings.TrimSpace(strings.Replace(string(s1), "\t", "  ", -1)), "\n",
	)
	l2 := strings.Split(
		strings.TrimSpace(strings.Replace(string(s2), "\t", "  ", -1)), "\n",
	)

	var n, w1, w2 int
	for i := 0; i < len(l1) && i < len(l2); i++ {
		if l1[i] == l2[i] {
			continue
		}
		if len(l1[i]) > w1 {
			w1 = len(l1[i])
		}
		if len(l2[i]) > w2 {
			w2 = len(l2[i])
		}
		if l := len(fmt.Sprintf("%d", i+1)); l > n {
			n = l
		}
	}
	for i := len(l2); i < len(l1); i++ {
		if len(l1[i]) > w1 {
			w1 = len(l1[i])
		}
		if l := len(fmt.Sprintf("%d", i+1)); l > n {
			n = l
		}
	}
	for i := len(l1); i < len(l2); i++ {
		if len(l2[i]) > w2 {
			w2 = len(l2[i])
		}
		if l := len(fmt.Sprintf("%d", i+1)); l > n {
			n = l
		}
	}

	var s []string
	for i := 0; i < len(l1) && i < len(l2); i++ {
		if l1[i] == l2[i] {
			continue
		}
		s = append(s, fmt.Sprintf(fmt.Sprintf("\t[%%%dd]!!! %%-%ds ! %%-%ds !!!", n, w1, w2), i+1, l1[i], l2[i]))
	}
	for i := len(l2); i < len(l1); i++ {
		s = append(s, fmt.Sprintf(fmt.Sprintf("\t[%%%dd]+++ %%-%ds + %%-%ds +++", n, w1, w2), i+1, l1[i], ""))
	}
	for i := len(l1); i < len(l2); i++ {
		s = append(s, fmt.Sprintf(fmt.Sprintf("\t[%%%dd]--- %%-%ds - %%-%ds ---", n, w1, w2), i+1, "", l2[i]))
	}

	return strings.Join(s, "\n")
}
