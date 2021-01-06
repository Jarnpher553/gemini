package csv

import (
	"bytes"
	"github.com/gocarina/gocsv"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

func Unmarshal(in []byte, out interface{}) error {
	return gocsv.UnmarshalBytes(in, out)
}

func Marshal(in interface{}, withExplain bool, withBom bool) ([]byte, error) {
	b, err := gocsv.MarshalBytes(in)
	if err != nil {
		return nil, err
	}

	if withExplain {
		v := reflect.ValueOf(in)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		t := v.Type().Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		bEx, err := getExplain(t, b, withBom)
		if err != nil {
			return nil, err
		}
		return bEx, nil
	}
	return b, nil
}

func getExplain(in reflect.Type, bIn []byte, withBom bool) ([]byte, error) {
	num := in.NumField()
	var out []byte
	for i := 0; i < num; i++ {
		name := in.Field(i).Tag.Get("name")
		out = []byte(string(out) + name + ",")
	}

	var buffer bytes.Buffer

	if withBom {
		_ = bom(&buffer)
	}

	buffer.Write(out[:len(out)-1])
	buffer.WriteByte('\r')
	buffer.WriteByte('\n')
	buffer.Write(bIn)

	return buffer.Bytes(), nil
}

func ClearExplain(reader io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	slice := strings.Split(string(b), "\r\n")
	out := strings.Join(slice[1:], "\r\n")
	return []byte(out), nil
}

func bom(buffer *bytes.Buffer) error {
	buffer.WriteByte(0xEF)
	buffer.WriteByte(0xBB)
	buffer.WriteByte(0xBF)

	return nil
}
