package model

import (
	"fmt"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/types"
	"github.com/go-redis/redis/v8"
)

type Course struct {
	Courseid   uint32 `gorm:"primaryKey;column:course_id"`
	Coursename string `gorm:"column:course_name"`
	Cap        uint32 `gorm:"column:cap"`
	Remaincap  uint32 `gorm:"column:remain_cap"`
	Teacherid  uint32 `gorm:"column:teacher_id"`
}

func (Course) TableName() string {
	return "course"
}

func (Course) Insert(req types.CreateCourseRequest) string {
	course := Course{
		Coursename: req.Name,
		Cap:        uint32(req.Cap),
		Remaincap:  uint32(req.Cap),
	}
	models.Db.Create(&course)
	Rdb.Set(Ctx, fmt.Sprintf("coursename%d", course.Courseid), course.Coursename, redis.KeepTTL)
	Rdb.Set(Ctx, fmt.Sprintf("courseteacher%d", course.Courseid), course.Teacherid, redis.KeepTTL)
	Rdb.Set(Ctx, fmt.Sprintf("course%d", course.Courseid), course.Remaincap, redis.KeepTTL)
	return fmt.Sprintf("%d", course.Courseid)
}

func (Course) CoursenameUsed(coursename string) bool {
	if err := models.Db.First(&Course{}, "course_name = ?", coursename).Error; err == nil { // find
		return true
	}
	return false
}
