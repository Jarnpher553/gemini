package main

const serviceTmpl = `
package services

import (
	"{{name}}/model"
	"{{name}}/middlewares"

	"github.com/Jarnpher553/gemini/copier"
	"github.com/Jarnpher553/gemini/erro"
	"github.com/Jarnpher553/gemini/model/dto"
	"github.com/Jarnpher553/gemini/repo"
	"github.com/Jarnpher553/gemini/service"
)

type {{title .Name}} struct {
	*service.BaseService
}

func (s *{{title .Name}}) Use(handler *service.Handler) {
	handler.UseMiddleware(service.AuthMiddleware(), middlewares.Permission("{{name}}.{{trimSuffix .Name "Service"}}"))
}

func (s *{{title .Name}}) Post(handler *service.Handler) service.HandlerFunc {
	return func(ctx *service.Ctx) {
		var in model.{{title .Dto.Request.Name}}

		if err := ctx.ShouldBind(&in); err != nil {
			ctx.Failure(erro.ErrReqContent, err)
			return
		}

		var entity model.{{title .Orm.Name}}
		_ = copier.Copy(&in, &entity)

		if err := s.Repo().Insert(&entity); err != nil {
			ctx.Failure(erro.ErrDbInsert, err)
			return
		}

		ctx.Success(nil)
	}
}

func (s *{{title .Name}}) Put(handler *service.Handler) service.HandlerFunc {
	return func(ctx *service.Ctx) {
		id := ctx.Param("id")

		var in model.{{title .Dto.Request.Name}}
		if err := ctx.ShouldBind(&in); err != nil {
			ctx.Failure(erro.ErrReqContent, err)
			return
		}

		if err := s.Repo().ModifyFunc(&model.{{title .Orm.Name}}{}, func(entity interface{}) {
			instant := entity.(*model.{{title .Orm.Name}})
			_ = copier.Copy(&in, instant)
		}, "id = ?", id); err != nil {
			ctx.Failure(erro.ErrDbModify, err)
			return
		}

		ctx.Success(nil)
	}
}

func (s *{{title .Name}}) Delete(handler *service.Handler) service.HandlerFunc {
	return func(ctx *service.Ctx) {
		id := ctx.Param("id")

		if err := s.Repo().Remove(&model.{{title .Orm.Name}}{}, "id = ?", id); err != nil {
			ctx.Failure(erro.ErrDbRemove, err)
			return
		}

		ctx.Success(nil)
	}
}

func (s *{{title .Name}}) GetList(handler *service.Handler) service.HandlerFunc {
	return func(ctx *service.Ctx) {
		var in dto.PagedIn
		if err := ctx.ShouldBind(&in); err != nil {
			ctx.Failure(erro.ErrReqContent, err)
			return
		}

		var entity []model.{{title .Dto.Response.Name}}
		var count int
		var err error
		if count, err = s.Repo().Query(&entity,
			true,
			repo.Model(&model.{{title .Orm.Name}}{}),
			repo.Page(in.PageNum, in.PerCount),
			repo.Order("created_time"),
		); err != nil {
			ctx.Failure(erro.ErrDbRead, err)
			return
		}

		ctx.Success(dto.PagedOut{
			TotalCount: count,
			PerCount:   in.PerCount,
			PageNum:    in.PageNum,
			QueryCount: len(entity),
			Rows:       entity,
		})
	}
}

func (s *{{title .Name}}) Get(handler *service.Handler) service.HandlerFunc {
	return func(ctx *service.Ctx) {
		id := ctx.Param("id")

		var out model.{{title .Dto.Response.Name}}
		if _, err := s.Repo().Query(&out, false, repo.Model(&model.{{title .Orm.Name}}{}), repo.Where("id = ?", id)); err != nil {
			ctx.Failure(erro.ErrDbRead, err)
			return
		}

		ctx.Success(&out)
	}
}
`
