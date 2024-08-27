package resp

// ResponseInfo is used to store simplified conversion details.
type ResponseInfo struct {
	Sensitivity float64
	Gain        float64
	Bias        float64
	Input       string
	Output      string
}

// Info builds a ResponseInfo for the given response base and sensor, database and stream collection details.
func (r *Resp) Info(responses ...string) (*ResponseInfo, error) {

	// simple response
	info := ResponseInfo{
		Sensitivity: 1.0,
		Gain:        1.0,
	}

	// run through the desired responses.
	for _, lookup := range responses {

		res, err := r.Type(lookup)
		if err != nil {
			return nil, err
		}

		// extract the response details into the response info.
		switch {
		case res.InstrumentSensitivity != nil:
			if info.Input == "" {
				info.Input = res.InstrumentSensitivity.InputUnits.Name
			}
			info.Output = res.InstrumentSensitivity.OutputUnits.Name
			info.Sensitivity *= res.InstrumentSensitivity.Value
		case res.InstrumentPolynomial != nil:
			if info.Input == "" {
				info.Input = res.InstrumentPolynomial.InputUnits.Name
			}
			info.Output = res.InstrumentPolynomial.OutputUnits.Name
			if coefs := res.PolynomialCoefficients(); len(coefs) == 2 {
				info.Bias = coefs[0].Value
				info.Gain = coefs[1].Value
			}
		}
	}

	return &info, nil
}
