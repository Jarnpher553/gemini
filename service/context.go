package service

import (
	"context"
	"github.com/Jarnpher553/gemini/erro"
	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/gemini/model/dto"
	"github.com/Jarnpher553/gemini/now"
	"github.com/Jarnpher553/gemini/uuid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	c.Data(http.StatusOK, "text/csv", data)
}

func (c *Ctx) Pdf(data []byte, filename string) {
	c.Header("Content-Disposition", "filename="+filename+".pdf")
	c.Data(http.StatusOK, "application/pdf", data)
}

func (c *Ctx) FileStream(data []byte, filename string) {
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (c *Ctx) Success(data interface{}) {
	c.response(erro.ErrSuccess, data, "", false)
}

func (c *Ctx) Failure(code int, err error, actual ...bool) {
	if len(actual) != 0 && actual[0] {
		c.response(code, nil, err.Error(), true)
	} else {
		c.response(code, nil, err.Error(), false)
	}
}

func (c *Ctx) Response(code int, data interface{}, err error) {
	c.response(code, data, err.Error(), false)
}

func (c *Ctx) response(code int, data interface{}, err string, actual bool) {
	var msg = erro.ErrMsg[code]
	if actual {
		msg = err
	}

	response := &dto.Response{
		ErrCode:   code,
		ErrMsg:    msg,
		Timestamp: now.New().Unix(),
		Data:      data,
	}

	log.Zap.Source(3).
		With(zap.Int("response.code", code)).
		With(zap.String("response.msg", erro.ErrMsg[code])).
		With(zap.String("response.err", err)).
		Info("response")

	c.JSON(http.StatusOK, response)
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

func (c *Ctx) User() interface{} {
	info := c.Request.Context().Value("auth_user")
	return info
}

func (c *Ctx) SetUser(user interface{}) {
	var cc context.Context
	cc = context.WithValue(c.Request.Context(), "auth_user", user)
	c.Request = c.Request.WithContext(cc)
}
