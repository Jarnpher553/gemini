package excel

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"reflect"
	"strconv"
)

func Gen(entities interface{}, titleColor ...string) (*bytes.Buffer, error) {
	var tc string
	if titleColor != nil && len(titleColor) != 0 {
		tc = titleColor[0]
	} else {
		tc = "#FFFF99"
	}
	tEn := reflect.TypeOf(entities)
	if tEn.Kind() != reflect.Slice {
		return nil, errors.New("input params type error")
	}
	vEn := reflect.ValueOf(entities)

	f := excelize.NewFile()
	t := vEn.Index(0).Type()
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
	styleID, _ := f.NewStyle(fmt.Sprintf(`{"font":{"bold":true},"fill":{"type":"pattern","color":["%s"],"pattern":1}}`, tc))
	endCol := string(65+count-1) + "1"
	f.SetCellStyle("Sheet1", "A1", endCol, styleID)
	for i := 0; i < vEn.Len(); i++ {
		var row []interface{}
		for j := 0; j < t.NumField(); j++ {
			tagValue := t.Field(j).Tag.Get("excel")
			if tagValue != "" {
				row = append(row, vEn.Index(i).Field(j).Interface())
			}
		}
		f.SetSheetRow("Sheet1", "A"+strconv.Itoa(i+2), &row)
	}

	buf, _ := f.WriteToBuffer()
	return buf, nil
}
