package expr

import (
	"errors"
	"go/constant"
	"go/token"
	"go/types"
)

var ErrEvalIsUnknown = errors.New("expression is unknown")
var ErrEvalNotBool = errors.New("expression does not return a bool expression")
var ErrEvalNotString = errors.New("expression does not return a string expression")

// parse the constant string via a token evaluation
func eval(s string) (constant.Value, error) {
	tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, s)
	if err != nil {
		return nil, err
	}
	if tv.Value.Kind() == constant.Unknown {
		return nil, ErrEvalIsUnknown
	}
	return tv.Value, nil
}

// ToFloat32 evaluates the given string as a float32, or it returns an error if the result is invalid or the expression is Unknown.
func ToFloat32(s string) (float32, error) {
	val, err := eval(s)
	if err != nil {
		return 0, err
	}
	res, _ := constant.Float32Val(val)
	return res, nil
}

// ToFloat64 evaluates the given string as a float64, or it returns an error if the result is invalid or the expression is Unknown.
func ToFloat64(s string) (float64, error) {
	val, err := eval(s)
	if err != nil {
		return 0, err
	}
	res, _ := constant.Float64Val(val)
	return res, nil
}

// ToInt64 evaluates the given string as an int64, or it returns an error if the result is invalid or the expression is Unknown.
func ToInt64(s string) (int64, error) {
	val, err := eval(s)
	if err != nil {
		return 0, err
	}
	res, _ := constant.Int64Val(val)
	return res, nil
}

// ToUint64 evaluates the given string as an uint64, or it returns an error if the result is invalid or the expression is Unknown.
func ToUint64(s string) (uint64, error) {
	val, err := eval(s)
	if err != nil {
		return 0, err
	}
	res, _ := constant.Uint64Val(val)
	return res, nil
}

// ToInt evaluates the given string as an int, or it returns an error if the result is invalid or the expression is Unknown.
func ToInt(s string) (int, error) {
	v, err := ToInt64(s)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// ToUint evaluates the given string as a uint, or it returns an error if the result is invalid or the expression is Unknown.
func ToUint(s string) (uint, error) {
	v, err := ToUint64(s)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// ToBool evaluates the given string as a bool, or it returns an error if the expression is Unknown.
func ToBool(s string) (bool, error) {
	val, err := eval(s)
	if err != nil {
		return false, err
	}
	if val.Kind() != constant.Bool {
		return false, ErrEvalNotBool
	}
	return constant.BoolVal(val), nil
}

// ToString evaluates the given string as a string, or it returns an error if the expression is Unknown.
func ToString(s string) (string, error) {
	val, err := eval(s)
	if err != nil {
		return "", err
	}
	if val.Kind() != constant.String {
		return "", ErrEvalNotString
	}
	return constant.StringVal(val), nil
}
