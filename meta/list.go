package meta

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type List interface {
	List()
	Sort() List
}

func MarshalList(l List) []byte {
	var b bytes.Buffer

	if lines, err := EncodeList(l.Sort()); err == nil {
		csv.NewWriter(&b).WriteAll(lines)
	}

	return b.Bytes()
}

func UnmarshalList(b []byte, l List) error {

	v, err := csv.NewReader(bytes.NewBuffer(b)).ReadAll()
	if err != nil {
		return err
	}

	return DecodeList(v, l)
}

func DecodeList(data [][]string, l List) error {
	rv := reflect.ValueOf(l)
	if rv.Kind() != reflect.Ptr {
		panic("list decode requires a pointer")
	}
	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Slice {
		panic("list decode requires a pointer to a slice")
	}

	// no data ...
	if !(len(data) > 1) {
		return nil
	}

	// where to start ...
	offset := rv.Len()

	// make space ...
	if rv.IsNil() {
		rv.Set(reflect.MakeSlice(rv.Type(), len(data)-1, len(data)-1))
	} else {
		rv.Set(reflect.AppendSlice(rv, reflect.MakeSlice(rv.Type(), len(data)-1, len(data)-1)))
	}

	// skip the header line ...
	for i, n := 1, len(data); i < n; i++ {
		// gather the slot to store the data
		ri := reflect.Indirect(rv.Index(offset + i - 1))
		if ri.Kind() != reflect.Struct {
			return fmt.Errorf("list decode requires a pointer to a slice of structs")
		}
		if ri.NumField() != len(data[i]) {
			return fmt.Errorf("list decode found incorrect number of fields")
		}
		// decode each field, in order ...
		for j := 0; j < ri.NumField(); j++ {
			v := ri.Field(j)
			s := strings.TrimSpace(data[i][j])

			if ri.Field(j).Type().AssignableTo(reflect.ValueOf(time.Time{}).Type()) {
				t, err := time.Parse(DateTimeFormat, s)
				if err != nil {
					return err
				}
				v.Set(reflect.ValueOf(t))
			} else {

				switch v.Kind() {
				case reflect.String:
					v.SetString(s)
				case reflect.Bool:
					b, err := strconv.ParseBool(s)
					if err != nil {
						return err
					}
					v.SetBool(b)
				case reflect.Int16:
					i, err := strconv.ParseInt(s, 10, 16)
					if err != nil {
						return err
					}
					v.SetInt(i)
				case reflect.Int32:
					i, err := strconv.ParseInt(s, 10, 32)
					if err != nil {
						return err
					}
					v.SetInt(i)
				case reflect.Float32:
					f, err := strconv.ParseFloat(s, 32)
					if err != nil {
						return err
					}
					v.SetFloat(f)
				case reflect.Float64:
					f, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return err
					}
					v.SetFloat(f)
				default:
					panic("invalid list field type: " + v.Kind().String())
				}
			}
		}
	}

	return nil
}

func EncodeList(l List) ([][]string, error) {
	var data [][]string

	rv := reflect.ValueOf(l)
	if rv.Kind() != reflect.Slice {
		panic("list encode requires a pointer to a slice")
	}

	if rv.IsNil() {
		return data, nil
	}

	for i, n := 0, rv.Len(); i < n; i++ {
		v := rv.Index(i).Interface()
		ri := reflect.ValueOf(v)
		if i == 0 {
			t := reflect.TypeOf(v)
			var header []string
			for j := 0; j < t.NumField(); j++ {
				f := t.Field(j)
				tags := strings.SplitAfter(strings.TrimSpace(t.Field(j).Tag.Get("csv")), ",")
				if len(tags) > 0 && len(tags[0]) > 0 {
					header = append(header, tags[0])
				} else {
					header = append(header, strings.Title(f.Name))
				}
			}

			data = append(data, header)
		}
		var line []string
		for j := 0; j < ri.NumField(); j++ {
			f := ri.Field(j)

			if ri.Field(j).Type().AssignableTo(reflect.ValueOf(time.Time{}).Type()) {
				line = append(line, f.Interface().(time.Time).UTC().Format(DateTimeFormat))
			} else {
				switch f.Kind() {
				case reflect.String:
					line = append(line, f.String())
				case reflect.Bool:
					line = append(line, strconv.FormatBool(f.Bool()))
				case reflect.Int16, reflect.Int32:
					line = append(line, strconv.FormatInt(f.Int(), 10))
				case reflect.Float64:
					line = append(line, strconv.FormatFloat(f.Float(), 'g', -1, 64))
				default:
					panic("invalid list field type: \"" + reflect.TypeOf(v).Field(j).Name + "\" [" + f.Kind().String() + "]")
				}
			}
		}
		data = append(data, line)
	}

	return data, nil
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

func LoadLists(dirname, filename string, l List) error {

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			if err := LoadList(path, l); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
