package models

import "gorm.io/gorm"

type CookieValue struct {
	UserID   uint64
	UserType int
}

//检查用户名和密码是否正确。
//用户名不存在或者用户名存在，但是密码不正确时。返回(false,nil,0)
//用户名和密码正确时，返回(true,nil,UserID)
//遇到非gorm.ErrRecordNotFound的错误时，返回(false,err,0)
func CheckAuth(username, password string) (bool, error, CookieValue) {
	var a CookieValue
	err := Db.Model(&Member{}).Where("username = ? AND password = ? AND is_active = ?",
		username, password, 1).First(&a).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err, CookieValue{}
	}
	if a.UserID > 0 {
		return true, nil, a
	}
	return false, nil, CookieValue{}
}
