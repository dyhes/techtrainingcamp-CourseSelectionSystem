package models

import (
	"fmt"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/logging"
	"gorm.io/gorm"
)

//课程Model
type Course struct {
	CourseID   uint64 `gorm:"primaryKey" json:"course_id"`
	CourseName string `json:"course_name"`
	Cap        uint   `json:"cap"`
	RemainCap  uint   `json:"remain_cap"`
	TeacherID  uint64 `json:"teacher_id"`
}

//根据课程名称查询课程是否已存在
func IsCourseExistByName(name string) (bool, error) {
	var course Course
	err := Db.Model(&Course{}).Select("course_id").
		Where("course_name = ?", name).First(&course).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		logging.Error(err)
		return false, err
	}
	if course.CourseID > 0 {
		return true, nil
	}

	return false, nil
}

//根据课程ID查询课程是否已存在
func IsCourseExistByID(id uint64) (bool, error) {
	var course Course
	err := Db.Model(&Course{}).Select("course_id").
		Where("course_id = ?", id).First(&course).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		logging.Error(err)
		return false, err
	}
	if course.CourseID > 0 {
		return true, nil
	}

	return false, nil
}

//根据课程ID查询课程是否已经被绑定
//课程不存在时，返回false,gorm.ErrRecordNotFound
//课程存在，但未绑定时，返回false,nil
//课程存在，但已绑定时，返回true,nil
func IsCourseBound(id uint64) (bool, error) {
	var course Course
	err := Db.Model(&Course{}).Select("teacher_id").
		Where("course_id = ?", id).First(&course).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logging.Error(err)
		}
		return false, err
	}
	if course.TeacherID > 0 {
		return true, nil
	}
	return false, nil
}

//创建课程并返回课程的ID
func CreateCourse(name string, cap uint) (uint64, error) {
	course := Course{CourseName: name, Cap: cap, RemainCap: cap}
	err := Db.Model(&Course{}).Create(&course).Error
	if err != nil {
		logging.Error(err)
		return 0, err
	}
	return course.CourseID, nil
}

type CourseInfo struct {
	CourseID   string `json:"CourseID"`
	CourseName string `json:"CourseName"`
	TeacherID  string `json:"TeacherID"`
}

//根据课程ID获取课程信息
func GetCourseInfo(id uint64) (CourseInfo, error) {
	var info Course
	err := Db.Model(&Course{}).Where("course_id = ?", id).First(&info).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logging.Error(err)
		}
		return CourseInfo{}, err
	}
	var res CourseInfo
	res.CourseID = fmt.Sprintf("%d", info.CourseID)
	res.CourseName = info.CourseName
	res.TeacherID = fmt.Sprintf("%d", info.TeacherID)
	return res, nil
}

//绑定课程。即将Course表中对应课程的teacher_id字段更更新为teacherID
func BindCourse(courseID, teacherID uint64) error {
	err := Db.Model(&Course{}).Where("course_id = ?", courseID).
		Update("teacher_id", teacherID).Error
	if err != nil {
		logging.Error(err)
		return err
	}
	return nil
}

//解绑课程。即将Course表对应课程的teacher_id字段更新为0
func UnBindCourse(courseID, teacherID uint64) error {
	var course Course
	err := Db.Model(&Course{}).Select("course_id").
		Where("course_id = ? AND teacher_id = ?", courseID, teacherID).First(&course).Error
	if err != nil {
		logging.Error(err)
		return err
	}

	err = Db.Model(&Course{}).Where("course_id = ? AND teacher_id = ?", courseID, teacherID).
		Update("teacher_id", 0).Error

	if err != nil {
		logging.Error(err)
		return err
	}
	return nil
}

//查询一个老师的全部课程
func GetTeacherCourses(teacherID uint64) ([]CourseInfo, error) {
	courses := make([]CourseInfo, 0)
	err := Db.Model(&Course{}).Where("teacher_id = ?", teacherID).Find(&courses).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(err)
		return nil, err
	}
	return courses, nil
}
