package cookie

import (
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/app"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

//若请求中没有名为camp-seesion的cookie，那么abort
func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	_, err := c.Cookie("camp-seesion")
	if err != nil {
		//err不为空，说明cookie不存在，此时abort，不能调用后面的路由
		appG.Response(http.StatusUnauthorized, e.LoginRequired, nil)
		c.Abort()
		return
	}
	c.Next()
}

//若当前登录的用户不是admin，那么直接abort
func Admin(c *gin.Context) {
	appG := app.Gin{C: c}
	cookie, _ := c.Cookie("camp-seesion")
	//cookie:{UserID+Username+UserType}
	userType, _ := strconv.Atoi(cookie[strings.LastIndex(cookie, "+")+1:])
	if models.UserType(userType) != models.Admin {
		appG.Response(http.StatusUnauthorized, e.PermDenied, nil)
		c.Abort()
		return
	}
	c.Next()
}
