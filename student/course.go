package student

import (
	"fmt"
	"net/http"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/model"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/types"
	"github.com/gin-gonic/gin"
)

// StudentNotExisted  ErrNo = 11 // 学生不存在

func Course(c *gin.Context) { // /api/v1/student/course
	var req types.GetStudentCourseRequest
	var res types.GetStudentCourseResponse
	c.BindQuery(&req)
	if model.Rdb.TTL(model.Ctx, fmt.Sprintf("studentcourse%s", req.StudentID)).Val().String()[:2] == "-2" {
		res.Code = types.StudentNotExisted
	} else {
		courses := model.Rdb.SMembers(model.Ctx, fmt.Sprintf("studentcourse%s", req.StudentID)).Val()
		for _, course := range courses {
			if course == "" {
				continue
			}
			res.Data.CourseList = append(res.Data.CourseList, types.TCourse{
				CourseID:  course,
				Name:      model.Rdb.Get(model.Ctx, fmt.Sprintf("coursename%s", course)).Val(),
				TeacherID: model.Rdb.Get(model.Ctx, fmt.Sprintf("courseteacher%s", course)).Val(),
			})
		}
	}
	c.JSON(http.StatusOK, res)
}

// type GetStudentCourseRequest struct {
// 	StudentID string `form:"StudentID" json:"StudentID"`
// }

// type GetStudentCourseResponse struct {
// 	Code ErrNo
// 	Data struct {
// 		CourseList []TCourse
// 	}
// }
