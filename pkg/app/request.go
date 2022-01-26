package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/e"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

//MakeErrors 向日志中写入表单验证时发生的错误
func MakeErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}

//根据cookie值判断当前登录的用户是不是admin
func IsAdmin(c *gin.Context) e.ErrNo {
	session, err := Store.Get(c.Request, "camp-seesion")
	if session.IsNew || err != nil {
		return e.LoginRequired
	}

	userType := session.Values["UserType"].(int)
	if userType != 1 {
		return e.PermDenied
	}

	return e.OK
}
