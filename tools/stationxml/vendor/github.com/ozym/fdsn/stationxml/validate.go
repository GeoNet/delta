package stationxml

type Validator interface {
	IsValid() error
}

func Validate(v Validator) error {
	return v.IsValid()
}

func ValidatePtr(v Validator) error {
	if v == nil {
		return nil
	}
	return v.IsValid()
}
