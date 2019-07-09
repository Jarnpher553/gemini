package now

import (
	"github.com/jinzhu/now"
	"time"
)

type Now struct {
	*now.Now
}

func New() *Now {
	return &Now{
		Now: now.New(time.Now()),
	}
}
