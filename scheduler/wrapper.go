package scheduler

func wrapper(sch *Scheduler, f func(*Options)) func() {
	return func() {
		f(sch.Options)
	}
}
