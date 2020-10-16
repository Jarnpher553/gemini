package main

import "fmt"

var versionTmpl = fmt.Sprintf(`
package services

import (
	"{{name}}/version"
	"github.com/Jarnpher553/gemini/service"
)

type VersionService struct {
	*service.BaseService
}

func (s *VersionService) GetBuild(handler *service.Handler) service.HandlerFunc {
	return func(ctx *service.Ctx) {
		ctx.Success(&struct {
			Release   string %s
			Commit    string %s
			BuildTime string %s
		}{
			version.Release,
			version.Commit,
			version.BuildTime,
		})
	}
}
`, "`json:\"release\"`", "`json:\"commit\"`", "`json:\"build_time\"`")
