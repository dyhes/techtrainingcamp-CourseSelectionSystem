package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/model"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/app"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/e"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/logging"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CreateCourseRequest struct {
	Name string `json:"Name" valid:"Required;MaxSize(255)"`
	Cap  uint   `json:"Cap"`
}

//@Summary 创建课程
//@Produce json
//@Param name query string false "Name"
//@Param cap query uint false "Cap"
//@Success 200 {string} json "{"code":200,"data":{course_id},"msg":{"ok"}}"
//@Router /api/v1/course/create [post]
func CreateCourse(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request CreateCourseRequest
	)
	//表单验证
	_, errCode := app.BindAndValid(c, &request, true)
	if errCode != e.OK {
		appG.Response(http.StatusOK, errCode, nil)
		return
	}

	//表单验证通过后，查看课程名称是否存在
	exist, err := models.IsCourseExistByName(request.Name)
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}
	if exist {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}

	//课程不存在，那么创建课程
	id, err := models.CreateCourse(request.Name, request.Cap)
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}
	model.Rdb.Set(model.Ctx, fmt.Sprintf("coursename%d", id), request.Name, redis.KeepTTL)
	model.Rdb.Set(model.Ctx, fmt.Sprintf("courseteacher%d", id), 0, redis.KeepTTL)
	model.Rdb.Set(model.Ctx, fmt.Sprintf("course%d", id), request.Cap, redis.KeepTTL)

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"CourseID": id})
}

type GetCourseRequest struct {
	CourseID string `form:"CourseID" valid:"Required"`
}

//@Summary 获取课程信息
//@Produce json
//@Param course_id query uint false "CourseList"
//@Success 200 {string} json "{"code":200,"data":{course},"msg":{"ok"}}"
//@Router /api/v1/course/get [get]
func GetCourse(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request GetCourseRequest
	)

	//表单验证
	_, errCode := app.BindAndValid(c, &request, false)
	if errCode != e.OK {
		appG.Response(http.StatusOK, errCode, nil)
		return
	}

	//表单验证成功后，获取课程信息
	id, err := strconv.ParseUint(request.CourseID, 10, 64)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusOK, e.ParamInvalid, nil)
		return
	}
	info, err := models.GetCourseInfo(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.CourseNotExisted, nil)
		} else {
			appG.Response(http.StatusOK, e.UnknownError, nil)
		}
		return
	}

	appG.Response(http.StatusOK, e.OK, info)
}

type BindCourseRequest struct {
	CourseID  string `json:"CourseID" valid:"Required"`
	TeacherID string `json:"TeacherID" valid:"Required"`
}

//@Summary  将教师与课程绑定
//@Produce json
//@Param course_id query uint64 false "CourseList"
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/teacher/bind_course [post]
func BindCourse(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request BindCourseRequest
	)

	//表单验证
	_, errCode := app.BindAndValid(c, &request, true)
	if errCode != e.OK {
		appG.Response(http.StatusOK, errCode, nil)
		return
	}

	//表单验证通过后，查询课程是否已被绑定
	courseID, err := strconv.ParseUint(request.CourseID, 10, 64)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusOK, e.ParamInvalid, nil)
		return
	}
	bound, _, err := models.IsCourseBound(courseID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.CourseNotExisted, nil)
		} else {
			appG.Response(http.StatusOK, e.UnknownError, nil)
		}
		return
	}
	if bound {
		appG.Response(http.StatusOK, e.CourseHasBound, nil)
		return
	}

	//若课程存在且未绑定时，绑定课程。需要注意的是，项目中的抢课部分是一个单独的算法题
	//因此，提供的teacher_id可能在数据库中根本不存在，故这里不做落库校验。
	teacherID, err := strconv.ParseUint(request.TeacherID, 10, 64)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusOK, e.ParamInvalid, nil)
		return
	}
	err = models.BindCourse(courseID, teacherID)
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}
	model.Rdb.Set(model.Ctx, fmt.Sprintf("courseteacher%d", courseID), teacherID, redis.KeepTTL)
	appG.Response(http.StatusOK, e.OK, nil)
}

//@Summary 将老师与课程解绑
//@Produce json
//@Param course_id query uint64 false "CourseList"
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/teacher/unbind_course [post]
func UnBindCourse(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request BindCourseRequest
	)

	//表单验证
	_, errCode := app.BindAndValid(c, &request, true)
	if errCode != e.OK {
		appG.Response(http.StatusOK, errCode, nil)
		return
	}

	//表单验证通过后，判断课程是否被绑定过
	courseID, err := strconv.ParseUint(request.CourseID, 10, 64)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusOK, e.ParamInvalid, nil)
		return
	}
	bound, teacherID, err := models.IsCourseBound(courseID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusOK, e.CourseNotExisted, nil)
		} else {
			appG.Response(http.StatusOK, e.UnknownError, nil)
		}
		return
	}
	if !bound {
		appG.Response(http.StatusOK, e.CourseNotBind, nil)
		return
	}

	//课程存在且被绑定时，尝试解绑课程
	teacherIDReq, err := strconv.ParseUint(request.TeacherID, 10, 64)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusOK, e.ParamInvalid, nil)
		return
	}
	if teacherID != teacherIDReq {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}

	err = models.UnBindCourse(courseID)
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}
	model.Rdb.Set(model.Ctx, fmt.Sprintf("courseteacher%d", courseID), 0, redis.KeepTTL)
	appG.Response(http.StatusOK, e.OK, nil)
}

type GetTeacherCourseRequest struct {
	TeacherID string `form:"TeacherID" valid:"Required"`
}

//@Summary 获取一个老师的全部课程
//@Produce json
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{courses},"msg":{"ok"}}"
//@Router /api/v1/teacher/get_course [GET]
func GetTeacherCourses(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request GetTeacherCourseRequest
	)

	//表单验证
	_, errCode := app.BindAndValid(c, &request, false)
	if errCode != e.OK {
		appG.Response(http.StatusOK, errCode, nil)
		return
	}

	//表单验证通过后，查询老师的全部课程
	teacherID, err := strconv.ParseUint(request.TeacherID, 10, 64)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusOK, e.ParamInvalid, nil)
		return
	}
	courses, err := models.GetTeacherCourses(teacherID)
	if err != nil {
		appG.Response(http.StatusOK, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"CourseList": courses})
}

type ScheduleCourseRequest struct {
	// key 为 teacherID , val 为老师期望绑定的课程 courseID 数组
	TeacherCourseRelationShip map[string][]string `json:"TeacherCourseRelationShip"`
}

func Schedule(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request = ScheduleCourseRequest{}
	)

	_, errCode := app.BindAndValid(c, &request, true)
	if errCode != e.OK {
		appG.Response(http.StatusOK, errCode, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, utility.MaxMatch(request.TeacherCourseRelationShip))
}
