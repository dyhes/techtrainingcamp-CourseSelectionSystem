package model

import (
	"fmt"
	"strconv"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"gorm.io/gorm"
)

type Pair struct {
	Studentid string
	Courseid  string
}

type Student struct {
	Studentid   uint32 `gorm:"primaryKey;column:student_id"`
	Studentname string `gorm:"column:student_name"`
}

func (Student) TableName() string {
	return "student"
}

func (Student) Insert(studentid uint32, studentname string) {
	models.Db.Create(Student{Studentid: studentid, Studentname: studentname})
	Rdb.SAdd(Ctx, fmt.Sprintf("studentcourse%d", studentid), "") // 为每位同学默认加入一个空的课程。
}

var BookCourseInfo = make(chan Pair, 10000)

func BookCourseHandler() {
	for {
		info := <-BookCourseInfo
		// studencourse表中加入该课程
		studentid, _ := strconv.Atoi(info.Studentid)
		courseid, _ := strconv.Atoi(info.Courseid)
		studentcourse := StudentCourse{
			Studentid:   uint32(studentid),
			Studentname: info.Studentid,
			Courseid:    uint32(courseid),
			Coursename:  info.Courseid,
		}
		models.Db.Create(studentcourse)
		// course 表中将该课程remaincap - 1
		models.Db.Model(&Course{}).Where("course_id = ?", info.Courseid).Update("remain_cap", gorm.Expr("remain_cap - ?", 1))
	}
}
