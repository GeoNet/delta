package sqlite

import (
	"fmt"
	"strings"
)

type QueryOpt func(int) (string, any)

func (q QueryOpt) K(i int) string {
	k, _ := q(i)
	return k
}

func (q QueryOpt) V() any {
	_, v := q(0)
	return v
}

func Option(k, v any) QueryOpt {
	return func(n int) (string, any) {
		return fmt.Sprintf(" %s = $%d", k, n+1), v
	}
}

func Code(v any) QueryOpt {
	return Option("Code", v)
}

func Network(v any) QueryOpt {
	return Option("Network", v)
}

func Station(v any) QueryOpt {
	return Option("Station", v)
}

func Location(v any) QueryOpt {
	return Option("Location", v)
}

func Mark(v any) QueryOpt {
	return Option("Mark", v)
}

func Make(v any) QueryOpt {
	return Option("Make", v)
}

func Model(v any) QueryOpt {
	return Option("Model", v)
}

func Serial(v any) QueryOpt {
	return Option("Serial", v)
}

func Sample(v any) QueryOpt {
	return Option("Sample", v)
}

func ParseBool(str string) (bool, bool) {
	switch strings.ToLower(str) {
	case "y", "yes", "true":
		return true, true
	case "n", "no", "false":
		return false, true
	default:
		return false, false
	}
}
