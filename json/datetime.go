package json

import (
	"time"
)

type Date struct {
	time.Time
}

func (date *Date) MarshalJSON() ([]byte, error) {
	if (*date).Time == (time.Time{}) {
		return []byte(`""`), nil
	}

	t := (*date).Time
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

func (date *Date) String() string {
	return (*date).Time.Format("2006-01-02")
}

func (date *Date) UnmarshalJSON(data []byte) (err error) {
	d := string(data[1 : len(data)-1])
	t, err := time.ParseInLocation("2006-01-02", string(d), time.Local)
	*date = Date{Time: t}
	return err
}
