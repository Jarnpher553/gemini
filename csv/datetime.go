package csv

import (
	"strings"
	"time"
)

type Date time.Time

func (date *Date) MarshalCSV() (string, error) {
	if *date == Date(time.Time{}) {
		return "", nil
	}

	return time.Time(*date).Format("2006/01/02"), nil
}

func (date *Date) String() string {
	return time.Time(*date).Format("2006/01/02")
}

func (date *Date) UnmarshalCSV(csv string) (err error) {
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
		*date = Date(t)
		return err
	}

	return nil
}

type DateTime time.Time

func (date *DateTime) MarshalCSV() (string, error) {
	if *date == DateTime(time.Time{}) {
		return "", nil
	}

	return time.Time(*date).Format("2006/01/02 15:04:05"), nil
}

func (date *DateTime) String() string {
	return time.Time(*date).Format("2006/01/02 15:04:05")
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	if csv != "" {
		csvSlice := strings.Split(strings.Split(csv, " ")[0], "/")

		if len(csvSlice[1]) == 1 {
			csvSlice[1] = "0" + csvSlice[1]
		}
		if len(csvSlice[2]) == 1 {
			csvSlice[2] = "0" + csvSlice[2]
		}

		csv = strings.Join(csvSlice, "/") + " " + strings.Split(csv, " ")[1]

		t, err := time.ParseInLocation("2006/01/02 15:04:05", csv, time.Local)
		*date = DateTime(t)
		return err
	}

	return nil
}
