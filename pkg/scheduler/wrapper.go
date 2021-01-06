package scheduler


func decorator(configuration *Configuration, f func(*Configuration)) func() {
	return func() {
		f(configuration)
	}
}
