package excel

import (
	"errors"
	"io"
)
import "github.com/360EntSecGroup-Skylar/excelize"

func Parse(file io.Reader, sheetName string) ([]map[string]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	content := f.GetRows(sheetName)

	if len(content) < 2 {
		return nil, errors.New("too few rows")
	}

	if len(content[0]) < 1 {
		return nil, errors.New("too few columns")
	}

	table := make([]map[string]string, 0)
	for i := 1; i < len(content); i++ {
		row := make(map[string]string, 0)
		for j := 0; j < len(content[i]); j++ {
			row[content[0][j]] = content[i][j]
		}
		table = append(table, row)
	}
	return table, nil
}
