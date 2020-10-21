package group

import (
	"github.com/Jarnpher553/gemini/log"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

var logger = log.Logger.Mark("group")

type G func() error

func Attach(f func() error) {
	g.Go(f)
}

func Run() {
	if err := g.Wait(); err != nil {
		logger.Fatal(err.Error())
	}
}
