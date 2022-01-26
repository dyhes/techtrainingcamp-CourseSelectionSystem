package model

type StudentCourse struct {
	Studentid   uint32 `gorm:"primary_key;column:student_id"`
	Studentname string `gorm:"column:student_name"`
	Courseid    uint32 `gorm:"column:course_id"`
	Coursename  string `gorm:"column:course_name"`
}

func (StudentCourse) TableName() string {
	return "student_course"
}
