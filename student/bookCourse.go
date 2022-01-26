package student

import (
	"fmt"
	"net/http"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/model"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/types"
	"github.com/gin-gonic/gin"
)

// StudentNotExisted  ErrNo = 11 // 学生不存在
// CourseNotExisted   ErrNo = 12 // 课程不存在
// StudentHasCourse   ErrNo = 14 // 学生有课程
// CourseNotAvailable ErrNo = 7  // 课程已满

func BookCourse(c *gin.Context) { // /api/v1/student/book_course
	var req types.BookCourseRequest
	var res types.BookCourseResponse
	c.BindJSON(&req)
	// fmt.Println(req.StudentID, req.CourseID, model.Rdb.TTL(model.Ctx, req.StudentID).Val().String()[:2])
	if model.Rdb.TTL(model.Ctx, fmt.Sprintf("studentcourse%s", req.StudentID)).Val().String()[:2] == "-2" { // 是否存在该学生
		res.Code = types.StudentNotExisted
	} else if model.Rdb.TTL(model.Ctx, fmt.Sprintf("course%s", req.CourseID)).Val().String()[:2] == "-2" { // 是否存在该课程
		res.Code = types.CourseNotExisted // 学生是否已经选择该课程
	} else if model.Rdb.SIsMember(model.Ctx, fmt.Sprintf("studentcourse%s", req.StudentID), req.CourseID).Val() {
		res.Code = types.StudentHasCourse
	} else { // 课程是否有容量
		remain := model.Rdb.Decr(model.Ctx, fmt.Sprintf("course%s", req.CourseID)).Val()
		if remain < 0 {
			model.Rdb.Incr(model.Ctx, fmt.Sprintf("course%s", req.CourseID))
			res.Code = types.CourseNotAvailable
		} else {
			model.Rdb.SAdd(model.Ctx, fmt.Sprintf("studentcourse%s", req.StudentID), req.CourseID)
			model.BookCourseInfo <- model.Pair{
				Studentid: req.StudentID,
				Courseid:  req.CourseID,
			}
		}
	}
	c.JSON(http.StatusOK, res)
}
