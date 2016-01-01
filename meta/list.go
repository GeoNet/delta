package meta

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
)

type ListEncoder interface {
	encode() [][]string
}
type ListDecoder interface {
	decode([][]string) error
}

type List interface {
	ListEncoder
	ListDecoder
}

func MarshalList(l ListEncoder) []byte {
	var b bytes.Buffer

	csv.NewWriter(&b).WriteAll(EncodeList(l))

	return b.Bytes()
}

func UnmarshalList(b []byte, l ListDecoder) error {

	v, err := csv.NewReader(bytes.NewBuffer(b)).ReadAll()
	if err != nil {
		return err
	}
	if err := DecodeList(v, l); err != nil {
		return err
	}

	return nil
}

func DecodeList(data [][]string, l ListDecoder) error {
	return l.decode(data)
}

func EncodeList(l ListEncoder) [][]string {
	return l.encode()
}

func LoadList(path string, l ListDecoder) error {

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

func StoreList(path string, l ListEncoder) error {

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
