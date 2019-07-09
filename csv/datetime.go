package csv

import (
	"strings"
	"time"
)

type DateTime time.Time

func (date *DateTime) MarshalCSV() (string, error) {
	if *date == DateTime(time.Time{}) {
		return "", nil
	}

	return time.Time(*date).Format("2006/01/02"), nil
}

func (date *DateTime) String() string {
	return time.Time(*date).Format("2006/01/02")
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	if csv != "" {
		csvSlice := strings.Split(csv, "/")

		if len(csvSlice[1]) == 1 {
			csvSlice[1] = "0" + csvSlice[1]
		}
		if len(csvSlice[2]) == 1 {
			csvSlice[2] = "0" + csvSlice[2]
		}

		csv = strings.Join(csvSlice, "/")

		t, err := time.ParseInLocation("2006/01/02", csv, time.Local)
		*date = DateTime(t)
		return err
	}

	return nil
}
