package routers

import (
	v1 "github.com/dyhes/techtrainingcamp-CourseSelectionSystem/routers/api/v1"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/student"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/teacher"
	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	//新建一个gin路由并绑定中间件
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())

	//设置路由
	apiv1 := g.Group("/api/v1")
	{
		apiv1.POST("/auth/login", v1.Login)   //登录
		apiv1.POST("/auth/logout", v1.Logout) //登出
		apiv1.GET("/auth/whoami", v1.WhoAmI)  //获取个人信息

		//成员
		apiv1.POST("/member/create", v1.CreateMember)
		apiv1.GET("/member/", v1.GetMember)
		apiv1.GET("/member/list", v1.GetMembers)
		apiv1.POST("/member/update", v1.UpdateMember)
		apiv1.POST("/member/delete", v1.DeleteMember)

		//排课
		apiv1.POST("/course/create", v1.CreateCourse)
		apiv1.GET("/course/get", v1.GetCourse)
		apiv1.POST("/teacher/bind_course", v1.BindCourse)
		apiv1.POST("/teacher/unbind_course", v1.UnBindCourse)
		apiv1.GET("/teacher/get_course", v1.GetTeacherCourses)
		apiv1.POST("/course/schedule", teacher.Schedule)

		//抢课
		apiv1.POST("/student/book_course", student.BookCourse)
		apiv1.GET("/student/course", student.Course)
	}

	return g
}
