package main

// Decoder is an interface to allow reading and then marshalling versions of stationxml
type Decoder interface {
	Decode(bool, []byte) ([]byte, error)
}

// Encode encodes the Root struct using the given Encoder.
func Decode(version string, verbose bool, data []byte) ([]byte, error) {

	var decoder Decoder
	switch version {
	case "1.0":
		decoder = Decoder10{}
	case "1.1":
		decoder = Decoder11{}
	case "1.2":
		decoder = Decoder12{}
	default:
		decoder = Decoder10{}
	}

	return decoder.Decode(verbose, data)
}
