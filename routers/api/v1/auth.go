package v1

import (
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/app"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"Username" valid:"Required;MinSize(8);MaxSize(20)"`
	Password string `json:"Password" valid:"Required;MinSize(8);MaxSize(20)"`
}

//@Summary 用户登录
//@Produce json
//@Param username query string false "UserName"
//@Param password query string false "Password"
//@Success 200 {string} json "{"code":200,"data":{UserID},"msg":{"ok"}}"
//@Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request LoginRequest
	)

	//表单验证
	_, errCode := app.BindAndValid(c, &request, true)
	if errCode != e.OK {
		appG.Response(http.StatusOK, errCode, nil)
		return
	}

	//表单验证通过，检查用户名和密码是否正确
	pass, err, value := models.CheckAuth(request.Username, request.Password)
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}
	if !pass {
		appG.Response(http.StatusOK, e.WrongPassword, nil)
		return
	}

	//登录成功，那么设置cookie
	session, _ := app.Store.Get(c.Request, "camp-seesion")
	session.Options.MaxAge = 3600 * 24 * 7 //一周
	session.Values["UserID"] = value.UserID
	session.Values["Username"] = request.Username
	session.Values["UserType"] = value.UserType
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}

	//登录成功时需要返回登录用户的ID
	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"UserID": value.UserID})
}

//Logout 登出
func Logout(c *gin.Context) {
	appG := app.Gin{C: c}

	//如果没有带上cookie，说明没有登录
	session, err := app.Store.Get(c.Request, "camp-seesion")
	if session.IsNew || err != nil {
		appG.Response(http.StatusOK, e.LoginRequired, nil)
		return
	}

	//删除cookie
	session.Options.MaxAge = -1
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, nil)
}

//WhoAmI 获取个人信息
func WhoAmI(c *gin.Context) {
	appG := app.Gin{C: c}

	//如果请求没有cookie，视为没有登录
	session, err := app.Store.Get(c.Request, "camp-seesion")
	if session.IsNew || err != nil {
		appG.Response(http.StatusOK, e.LoginRequired, nil)
		return
	}

	//根据UserID查询个人信息
	userInfo, err := models.GetUserInfo(session.Values["UserID"].(uint64))
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, userInfo)
}
