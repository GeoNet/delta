// Package expr provides helpful conversion of mathematic expressions to base types.
//
// The goal of this package is to behave in a similar manner to the strconv package
// for numerical type conversion.
//
// It allows for code like this:
//
//	v, err := ToFloat64("1.0 + 1.0")
//	// v == 2.0
package expr
