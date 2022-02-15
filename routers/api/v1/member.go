package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/model"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/app"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/e"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/logging"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/utility"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateMemberRequest struct {
	Nickname string `json:"Nickname" valid:"Required;MinSize(4);MaxSize(20)"`
	Username string `json:"Username" valid:"Required;Match(/^[A-Za-z]{8,20}$/)"`
	Password string `json:"Password" valid:"Required;MinSize(8);MaxSize(20);PasswordCheck"`
	UserType int    `json:"UserType" valid:"Required;Range(1,3)"`
}

//@Summary 创建成员。只有管理员能够访问，请求中需要带上cookie
//@Produce json
//@Param username query string false "UserName"
//@Param password query string false "Password"
//@Param nickname query string false "Nickname"
//@Param user_type query int false "UserType"
//@Success 200 {string} json "{"code":200,"data":{"user_id"},"msg":{"ok"}}"
//@Router /api/v1/member/create [post]
func CreateMember(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request CreateMemberRequest
	)

	//权限验证
	if errCode := app.IsAdmin(c); errCode != e.OK {
		appG.Response(http.StatusUnauthorized, e.PermDenied, nil)
		return
	}

	//表单验证
	funcs := app.CustomFunc{
		"PasswordCheck": utility.PasswordCheck,
	}
	httpCode, errCode := app.BindAndValidCustom(c, &request, funcs, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过后，检查username是否已经存在
	exist, err := models.IsUserExistByName(request.Username)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}
	if exist {
		appG.Response(http.StatusOK, e.UserHasExisted, nil)
		return
	}

	//用户名不存在，则创建新用户
	user := models.Member{
		Username: request.Username,
		Password: request.Password,
		Nickname: request.Nickname,
		UserType: request.UserType,
	}
	if err = models.CreateMember(&user); err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}
	if request.UserType == 2 {
		model.Rdb.SAdd(model.Ctx, fmt.Sprintf("studentcourse%d", user.UserID), "")
	}
	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"UserID:": user.UserID})
}

type UpdateMemberRequest struct {
	UserID   string `json:"UserID"  valid:"Required"`
	Nickname string `json:"Nickname" valid:"Required;MinSize(4);MaxSize(20)"`
}

//@Summary 更新用户的昵称
//@Produce json
//@Param user_id query uint false "UserID"
//@Param nickname query string false "Nickname"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router	/api/v1/member/update [post]
func UpdateMember(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request UpdateMemberRequest
	)
	//表单验证
	httpCode, errCode := app.BindAndValid(c, &request, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功后，检查username是否存在以及是否已删除
	id, err := strconv.ParseUint(request.UserID, 10, 64)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ParamInvalid, nil)
		return
	}
	deleted, err := models.IsUserDeleted(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.UserNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}
	if deleted {
		appG.Response(http.StatusOK, e.UserHasDeleted, nil)
		return
	}

	//username存在且没有被删除，那么修改用户名
	err = models.UpdateMember(id, request.Nickname)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, nil)
}

type DelMemberRequest struct {
	UserID string `json:"UserID" valid:"Required"`
}

//@Summary  删除用户(软删除)
//@Produce json
//@Param user_id query uint false "UserID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/member/delete [post]
func DeleteMember(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request DelMemberRequest
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &request, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功，那么查询用户是否存在以及是否是已删除状态
	id, err := strconv.ParseUint(request.UserID, 10, 64)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusBadRequest, e.ParamInvalid, nil)
		return
	}
	deleted, err := models.IsUserDeleted(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.UserNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}
	if deleted {
		appG.Response(http.StatusOK, e.UserHasDeleted, nil)
		return
	}

	//用户存在且是未删除状态时，删除用户
	err = models.DeleteMember(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, nil)
}

type GetMemberRequest struct {
	UserID string `form:"UserID" valid:"Required"`
}

//@Summary	获取单个成员信息
//@Produce json
//@Param user_id query uint false "UserID"
//@Success 200 {string} json "{"code":200,"data":{userinfo},"msg":{"ok"}}"
//@Router	/api/v1/member/ [get]
func GetMember(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request GetMemberRequest
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &request, false)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功后，查询用户名是否存在
	id, err := strconv.ParseUint(request.UserID, 10, 64)
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusBadRequest, e.ParamInvalid, nil)
		return
	}
	deleted, err := models.IsUserDeleted(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.UserNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}
	if deleted {
		appG.Response(http.StatusBadRequest, e.UserHasDeleted, nil)
		return
	}

	//用户名存在且没有被删除，那么获取用户信息
	info, err := models.GetUserInfo(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, info)
}

type GetMemberListRequest struct {
	Offset int `valid:"Min(0)"`
	Limit  int `valid:"Min(-1)"`
}

//@Summary	批量获取成员信息
//@Produce json
//@Param offset query uint false "OffSet"
//@Param limit query uint false "Limit"
//@Success 200 {string} json "{"code":200,"data":{userinfo},"msg":{"ok"}}"
//@Router /api/v1/member/list [get]
func GetMembers(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		request GetMemberListRequest
	)

	//表单验证
	_, b1 := c.GetQuery("Offset")
	_, b2 := c.GetQuery("Limit")
	if !b1 || !b2 {
		appG.Response(http.StatusBadRequest, e.ParamInvalid, nil)
		return
	}
	httpCode, errCode := app.BindAndValid(c, &request, false)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//批量获取信息
	members, err := models.GetUserInfoList(request.Offset, request.Limit)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"MemberList": members})
}
