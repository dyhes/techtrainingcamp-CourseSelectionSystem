package utility

import (
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"strconv"
	"strings"
)

//cookie的值为：UserID+Username+UserType

func GetUserIDFromCookie(value string) uint64 {
	//找到第一个+的位置
	index := strings.Index(value, "+")
	idStr := value[:index]
	id, _ := strconv.ParseUint(idStr, 10, 64)
	return id
}

func GetUsernameFromCookie(value string) string {
	//找到第一个和最后一个+的位置
	first := strings.Index(value, "+")
	last := strings.LastIndex(value, "+")
	return value[first+1 : last]
}

func GetUserTypeFromCookie(value string) models.UserType {
	//找到最后一个+所在的位置
	index := strings.LastIndex(value, "+")
	userType, _ := strconv.Atoi(value[index+1:])
	var res models.UserType
	switch models.UserType(userType) {
	case models.Admin:
		res = models.Admin
	case models.Student:
		res = models.Student
	case models.Teacher:
		res = models.Teacher
	}
	return res
}
