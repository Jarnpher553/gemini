package ulid

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func New() string {
	t := time.Now()
	return ulid.MustNew(ulid.Timestamp(t), ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)).String()
}
