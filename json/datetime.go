package json

import (
	"time"
)

type Date time.Time

func (date *Date) MarshalJSON() ([]byte, error) {
	if *date == Date(time.Time{}) {
		return []byte(`""`), nil
	}

	t := time.Time(*date)
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

func (date *Date) String() string {
	return time.Time(*date).Format("2006-01-02")
}

func (date *Date) UnmarshalJSON(data []byte) (err error) {
	d := string(data[1 : len(data)-1])
	if d == "" {
		return nil
	}
	t, err := time.ParseInLocation("2006-01-02", string(d), time.Local)
	*date = Date(t)
	return err
}

type DateTime time.Time

func (date *DateTime) MarshalJSON() ([]byte, error) {
	if *date == DateTime(time.Time{}) {
		return []byte(`""`), nil
	}

	t := time.Time(*date)
	return []byte(`"` + t.Format("2006-01-02 15:04:05") + `"`), nil
}

func (date *DateTime) String() string {
	return time.Time(*date).Format("2006-01-02 15:04:05")
}

func (date *DateTime) UnmarshalJSON(data []byte) (err error) {
	d := string(data[1 : len(data)-1])
	if d == "" {
		return nil
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", string(d), time.Local)
	*date = DateTime(t)
	return err
}
