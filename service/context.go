package service

import (
	"context"
	"fmt"
	"github.com/Jarnpher553/micro-core/erro"
	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/micro-core/model/dto"
	"github.com/Jarnpher553/micro-core/now"
	"github.com/Jarnpher553/micro-core/uuid"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Ctx struct {
	*gin.Context
}

func Wrapper(f HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Ctx{c}
		f(ctx)
	}
}

type HandlerFunc func(*Ctx)

func (c *Ctx) Csv(data []byte, filename string) {
	c.Header("Content-Disposition", "attachment; filename="+filename+".csv")
	c.Data(erro.ErrSuccess, "text/csv", data)
}

func (c *Ctx) Success(data interface{}) {
	log.Logger.Caller(2).Infoln(erro.ErrSuccess, erro.ErrMsg[erro.ErrSuccess], fmt.Sprintf("%#v", data))
	c.response(erro.ErrSuccess, data)
}

func (c *Ctx) Failure(code int, err error) {
	log.Logger.Caller(2).Errorln(code, erro.ErrMsg[code], err)
	c.response(code, nil)
}

func (c *Ctx) Response(code int, data interface{}) {
	c.response(code, data)
}

func (c *Ctx) response(code int, data interface{}) {
	c.JSON(http.StatusOK, &dto.Response{
		ErrCode:   code,
		ErrMsg:    erro.ErrMsg[code],
		Timestamp: now.New().Unix(),
		Data:      data,
	})
}

func (c *Ctx) UserGUID() (uuid.GUID, bool) {
	id, ok := c.Request.Context().Value("auth_user_guid").(uuid.GUID)
	return id, ok
}

func (c *Ctx) SetUserGUID(guid uuid.GUID) {
	var cc context.Context
	cc = context.WithValue(c.Request.Context(), "auth_user_guid", guid)
	c.Request = c.Request.WithContext(cc)
}

func (c *Ctx) UserID() (int, bool) {
	id, ok := c.Request.Context().Value("auth_user_id").(int)
	return id, ok
}

func (c *Ctx) SetUserID(id int) {
	var cc context.Context
	cc = context.WithValue(c.Request.Context(), "auth_user_id", id)
	c.Request = c.Request.WithContext(cc)
}

func (c *Ctx) UserInfo() (string, bool) {
	info, ok := c.Request.Context().Value("auth_user_info").(string)
	return info, ok
}

func (c *Ctx) SetUserInfo(userInfo string) {
	var cc context.Context
	cc = context.WithValue(c.Request.Context(), "auth_user_info", userInfo)
	c.Request = c.Request.WithContext(cc)
}
