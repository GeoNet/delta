package meta

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
)

type List interface {
	encode() [][]string
	decode([][]string) error
}

func MarshalList(l List) []byte {
	var b bytes.Buffer

	csv.NewWriter(&b).WriteAll(EncodeList(l))

	return b.Bytes()
}

func UnmarshalList(b []byte, l List) error {

	v, err := csv.NewReader(bytes.NewBuffer(b)).ReadAll()
	if err != nil {
		return err
	}
	if err := DecodeList(v, l); err != nil {
		return err
	}

	return nil
}

func DecodeList(data [][]string, l List) error {
	return l.decode(data)
}

func EncodeList(l List) [][]string {
	return l.encode()
}

func LoadList(path string, l List) error {

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	if err := DecodeList(data, l); err != nil {
		return err
	}

	return nil
}

func StoreList(path string, l List) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = csv.NewWriter(file).WriteAll(EncodeList(l))
	if err != nil {
		return err
	}

	return nil
}
