package retry

import (
	"time"
)

type Worker struct {
	Work     func() error
	Retry    int
	TimeSpan time.Duration
}

func (rw *Worker) Run() {
	retry := 0
	for {
		retry++
		err := rw.Work()
		if err == nil {
			break
		}

		if retry == rw.Retry {
			break
		}
		<-time.After(rw.TimeSpan)
	}
}

type WorkerBack struct {
	Worker
	Work func() (interface{}, error)
}

func (rw *WorkerBack) Run() (interface{}, error) {
	retry := 0
	var ret interface{}
	var err error
	for {
		retry++
		ret, err = rw.Work()
		if err == nil {
			break
		}

		if retry == rw.Retry {
			break
		}
		<-time.After(rw.TimeSpan)
	}
	return ret, err
}
