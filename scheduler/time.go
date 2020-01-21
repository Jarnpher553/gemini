package scheduler

import (
	"github.com/Jarnpher553/gron/xtime"
	"time"
)

const (
	//Second has 1 * 1e9 nanoseconds
	Second time.Duration = xtime.Second
	//Minute has 60 seconds
	Minute time.Duration = xtime.Minute
	//Hour has 60 minutes
	Hour time.Duration = xtime.Hour
	//Day has 24 hours
	Day time.Duration = xtime.Day
	//Week has 7 days
	Week time.Duration = xtime.Week
)
