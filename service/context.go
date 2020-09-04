package service

import (
	"context"
	"fmt"
	"github.com/Jarnpher553/gemini/erro"
	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/gemini/model/dto"
	"github.com/Jarnpher553/gemini/now"
	"github.com/Jarnpher553/gemini/uuid"
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
	dataStr := fmt.Sprintf("%#v", data)
	if len(dataStr) > 255 {
		dataStr = dataStr[:255]
	}
	log.Zap.Caller(2).Info(log.Message(erro.ErrSuccess, erro.ErrMsg[erro.ErrSuccess], dataStr+"..."))
	c.response(erro.ErrSuccess, data)
}

func (c *Ctx) Failure(code int, err error, actual ...bool) {
	log.Zap.Caller(2).Error(log.Message(code, erro.ErrMsg[code], err))

	if len(actual) != 0 && actual[0] {
		c.response(code, nil, err.Error())
	} else {
		c.response(code, nil)
	}
}

func (c *Ctx) Response(code int, data interface{}) {
	log.Zap.Caller(2).Error(log.Message(code, erro.ErrMsg[code]))

	c.response(code, data)
}

func (c *Ctx) response(code int, data interface{}, err ...string) {
	var e = erro.ErrMsg[code]
	if len(err) != 0 && err[0] != "" {
		e = err[0]
	}
	c.JSON(http.StatusOK, &dto.Response{
		ErrCode:   code,
		ErrMsg:    e,
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

func (c *Ctx) User() interface{} {
	info := c.Request.Context().Value("auth_user")
	return info
}

func (c *Ctx) SetUser(user interface{}) {
	var cc context.Context
	cc = context.WithValue(c.Request.Context(), "auth_user", user)
	c.Request = c.Request.WithContext(cc)
}

func (c *Ctx) Bind(obj interface{}) error {
	err := c.ShouldBind(obj)

	if err != nil {
		return err
	}

	dataStr := fmt.Sprintf("%#v", obj)
	if len(dataStr) > 255 {
		dataStr = dataStr[:255]
	}

	log.Zap.Caller(2).Info(dataStr + "...")
	return err
}
