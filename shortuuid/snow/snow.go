package snow

import (
	"fmt"
	"github.com/sony/sonyflake"
	"time"
)

var Flake = sonyflake.NewSonyflake(sonyflake.Settings{
	StartTime: time.Now(),
})

func NextID() string {
	id, _ := Flake.NextID()
	idStr := fmt.Sprintf("%d", id)
	return idStr
}
