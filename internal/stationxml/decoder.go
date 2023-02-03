package stationxml

// Encode encodes the Root struct using the given Encoder.
func Decode(version string, data []byte) ([]byte, error) {
	switch version {
	case "1.0":
		return Decode10(data)
	case "1.1":
		return Decode11(data)
	case "1.2":
		return Decode12(data)
	default:
		return Decode10(data)
	}
}
