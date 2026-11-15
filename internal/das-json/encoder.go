package dasjson

import (
	"io"
	"os"
)

// Encode encodes the Root struct using the given Encoder.
func (r Root) MarshalVersion(version string, schema string) ([]byte, error) {
	switch version {
	case "2.0":
		return Encode20(r, schema)
	default:
		return Encode20(r, schema)
	}
}

// Write encodes and writes the output from the given Encoder for a Root struct.
func (r Root) Write(wr io.Writer, version string, schema string) error {
	res, err := r.MarshalVersion(version, schema)
	if err != nil {
		return err
	}
	if _, err := wr.Write(res); err != nil {
		return err
	}

	return nil
}

// WriteFile encodes and stores the output from the given Encoder for a Root struct.
func (r Root) WriteFile(path string, version string, schema string) error {
	file, err := os.Create(path) //nolint:gosec // disable G304
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close() // ignore close error as it will likely be related to an existing error
	}()

	if err := r.Write(file, version, schema); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
