package excel

import (
	"bytes"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
	"reflect"
	"strconv"
)

func Gen(entities interface{}) (*bytes.Buffer, error) {
	tEn := reflect.TypeOf(entities)
	if tEn.Kind() != reflect.Slice {
		return nil, errors.New("input params type error")
	}
	vEn := reflect.ValueOf(entities)

	f := excelize.NewFile()
	t := vEn.Index(0).Type()
	v := vEn.Index(0)
	var title []string

	var count int
	for i := 0; i < t.NumField(); i++ {
		tagValue := t.Field(i).Tag.Get("excel")
		if tagValue != "" {
			title = append(title, tagValue)
			count++
		}
	}

	f.SetSheetRow("Sheet1", "A1", &title)
	styleID, _ := f.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","color":["#E0EBF5"],"pattern":1}}`)
	endCol := string(65+count-1) + "1"
	f.SetCellStyle("Sheet1", "A1", endCol, styleID)
	for i := 0; i < vEn.Len(); i++ {
		var row []interface{}
		for j := 0; j < t.NumField(); j++ {
			tagValue := t.Field(i).Tag.Get("excel")
			if tagValue != "" {
				row = append(row, v.Field(j).Interface())
			}
		}
		f.SetSheetRow("Sheet1", "A"+strconv.Itoa(i+2), &row)
	}

	buf, _ := f.WriteToBuffer()
	return buf, nil
}
