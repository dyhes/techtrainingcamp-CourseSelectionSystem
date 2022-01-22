package routers

import (
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/setting"
	v1 "github.com/dyhes/techtrainingcamp-CourseSelectionSystem/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	//新建一个gin路由并绑定中间件
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())
	gin.SetMode(setting.RunMode)

	//设置路由
	apiv1 := g.Group("/api/v1/")
	{
		//登录
		apiv1.POST("/auth/login", v1.Login)
		apiv1.POST("/auth/logout", v1.Logout)
		apiv1.GET("/auth/whoami", v1.WhoAmI)
	}
	return g
}
