package models

func GetCourseRemainCap(courseID uint64) (uint, error) {
	var c Course
	err := Db.Model(&Course{}).Select("remain_cap").
		Where("course_id = ?", courseID).First(&c).Error
	return c.RemainCap, err
}

///*
//定时将缓存数据存入数据库
//*/
//
//var (
//	ModifiedCourse  utility.Queue
//	ModifiedStudent utility.Queue
//	mutex           sync.Mutex
//)
//
//type StudentCourse struct {
//
//}
//
////更新课程的剩余名额
//func UpdateCourseRemainCap(courseID uint64, cap uint) error {
//	return db.Model(&Course{}).Where("course_id = ?", courseID).
//		Update("remain_cap", cap).Error
//}
//
//func UpdateStudentCourse(studentID uint64, courseList string) error {
//	return db.Model(&StudentCourse{}).Where("student_id = ?", studentID).
//		Update("course_list", courseList).Error
//}
//
////应用缓存更改到数据库
//func ApplyModifyToDataBase() {
//	mutex.Lock()
//	defer mutex.Unlock()
//
//	//更新课程的剩余名额
//	for !ModifiedCourse.Empty() {
//		courseID := (*ModifiedCourse.Pop()).(uint64)
//		v, _ := CourseCache.cache[courseID]
//		err := UpdateCourseRemainCap(courseID, v.remainCap)
//		if err != nil {
//			logging.Error(err)
//		}
//	}
//	//更新学生课表
//	for !ModifiedStudent.Empty() {
//		studentID := (*ModifiedStudent.Pop()).(uint64)
//		cSet, _ := StudentCourseCache.cache[studentID]
//		err := UpdateStudentCourse(studentID, cSet.GetCourseList(";"))
//		if err != nil {
//			logging.Error(err)
//		}
//	}
//}
//
////清空本地缓存
//func ClearCache() {
//	mutex.Lock()
//	defer mutex.Unlock()
//	ModifiedCourse = utility.Queue{}
//	ModifiedStudent = utility.Queue{}
//	StudentCourseCache = studentCourseCache{cache: map[uint64]*courseSet{}}
//	CourseCache = courseCache{cache: map[uint64]*courseCap{}}
//}
