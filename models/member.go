package models

import (
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/logging"
	"gorm.io/gorm"
)

//成员类型
type UserType int

const (
	Admin UserType = iota + 1
	Student
	Teacher
)

//成员Model
type Member struct {
	UserID   uint64 `gorm:"primaryKey" json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	UserType int    `json:"user_type"`
	IsActive int    `json:"is_active" gorm:"default:1"`
}

//用于返回查询到的成员信息
type MemberInfo struct {
	UserID   uint64 //`json:"user_id"`
	Username string //`json:"username"`
	Nickname string //`json:"nickname"`
	UserType int
}

//学生Model
type student struct {
	StudentID   uint64
	StudentName string
}

//教师Model
type teacher struct {
	TeacherID   uint64
	TeacherName string
}

//通过ID查询用户是否已存在。
//若存在，返回(true,nil)；若不存在，返回(false,nil)；若发生错误，返回(false,err)
func IsUserExistByID(id uint64, userType int) (bool, error) {
	type temp struct {
		UserID uint64
	}
	var t temp
	err := Db.Model(&Member{}).Where("user_id = ? AND user_type = ? AND is_active = ?",
		id, userType, 1).First(&t).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if t.UserID > 0 {
		return true, nil
	}
	return false, nil
}

//通过username查询用户是否已存在
//若存在，返回(true,nil)；若不存在，返回(false,nil)；若发生错误，返回(false,err)
func IsUserExistByName(username string) (bool, error) {
	type temp struct {
		UserID uint64
	}
	var t temp
	err := Db.Model(&Member{}).Where("username = ?", username).First(&t).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if t.UserID > 0 {
		return true, nil
	}
	return false, nil
}

//判断成员是否被删除。
//若用户不存在，返回(false,gorm.ErrRecordNotFound);
//若发生错误，返回(false,err)
//若用户存在，但是已删除状态(is_active=0)，返回(true,nil);
//若用户存在，但不是删除状态(is_active=1)，返回(false,nil)
func IsUserDeleted(id uint64) (bool, error) {
	type temp struct {
		IsActive int
	}
	var t temp
	err := Db.Model(&Member{}).Where("user_id = ?", id).First(&t).Error
	if err != nil {
		return false, err
	}
	if t.IsActive == 0 {
		return true, nil
	}
	return false, nil
}

//创建成员
func CreateMember(user *Member) error {
	//创建成员
	err := Db.Model(&Member{}).Create(user).Error
	if err != nil {
		return err
	}

	//根据成员类型再创建学生或者老师
	switch UserType(user.UserType) {
	case Student:
		s := student{StudentID: user.UserID, StudentName: user.Username}
		err = Db.Create(&s).Error
	case Teacher:
		t := teacher{TeacherID: user.UserID, TeacherName: user.Username}
		err = Db.Create(&t).Error
	}

	if err != nil {
		return err
	}
	return nil
}

//更新成员的nickname
func UpdateMember(id uint64, nickname string) error {
	err := Db.Model(&Member{}).Where("user_id = ?", id).Update("nickname", nickname).Error
	if err != nil {
		return err
	}
	return nil
}

//删除成员(软删除)
func DeleteMember(id uint64) error {
	err := Db.Model(&Member{}).Where("user_id = ?", id).Update("is_active", 0).Error
	if err != nil {
		return err
	}
	return nil
}

//获取单个成员的信息
func GetUserInfo(id uint64) (MemberInfo, error) {
	var u MemberInfo
	err := Db.Model(&Member{}).Where("user_id = ?", id).First(&u).Error
	if err != nil {
		return MemberInfo{}, err
	}
	return u, nil
}

//批量返回成员信息
func GetUserInfoList(offset int, limit int) (members []MemberInfo, err error) {
	if limit == 0 {
		return []MemberInfo{}, nil
	}
	err = Db.Model(&Member{}).Where("is_active = ?", 1).Offset(offset).
		Limit(limit).Find(&members).Error
	if err != nil {
		logging.Error(err)
		return members, err
	}
	return members, nil
}
