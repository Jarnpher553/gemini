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
	t, err := time.ParseInLocation("2006-01-02", string(d), time.Local)
	*date = Date(t)
	return err
}
