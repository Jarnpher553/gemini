[![Go Report Card](https://goreportcard.com/badge/github.com/Jarnpher553/gemini)](https://goreportcard.com/report/github.com/Jarnpher553/gemini)
![GitHub](https://img.shields.io/badge/license-MIT-brightgreen)

# gemini

## Introduction

GEMINI is a microservice framework

## Quick Start

```go
package main

import (
	"geminiv2demo/internal/services/user"
	"github.com/Jarnpher553/gemini/pkg/config"
	"github.com/Jarnpher553/gemini/pkg/router"
	"github.com/Jarnpher553/gemini/pkg/server"
	"github.com/Jarnpher553/gemini/pkg/service"
)

func main() {
	config.Args()
	c := config.Conf()

	restServer := server.Default(server.Env(c.GetString("server.env")), server.Addr(c.GetString("server.addr")), server.Route(func() *router.Router {
		r := router.New(router.Root("api"), router.Area())
		r.Assign(service.NewService(&user.DemoService{})).
			Register()

		return r
	}))
	restServer.Run()
}
```
