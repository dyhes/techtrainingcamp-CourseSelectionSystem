package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/e"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CustomFunc map[string]func(v *validation.Validation, obj interface{}, key string)

//表单验证。返回http状态码以及参数检验的错误代码。
func BindAndValid(c *gin.Context, obj interface{}, json bool) (int, e.ErrNo) {
	//将请求中的对应参数绑定到了form中
	var err error
	if json {
		err = c.BindJSON(obj)
	} else {
		err = c.ShouldBindQuery(obj)
	}

	if err != nil {
		return http.StatusBadRequest, e.ParamInvalid
	}

	valid := validation.Validation{}
	right, err := valid.Valid(obj)
	if err != nil {
		logging.Error(err)
		return http.StatusInternalServerError, e.UnknownError
	}
	if !right {
		MakeErrors(valid.Errors)
		return http.StatusBadRequest, e.ParamInvalid
	}

	return http.StatusOK, e.OK
}

//支持自定义表单验证函数
func BindAndValidCustom(c *gin.Context, obj interface{}, funcs CustomFunc, json bool) (int, e.ErrNo) {
	//先将请求中的参数绑定到form中
	var err error
	if json {
		err = c.BindJSON(obj)
	} else {
		err = c.ShouldBindQuery(obj)
	}
	if err != nil {
		return http.StatusBadRequest, e.ParamInvalid
	}

	//设置自定义的表单验证函数
	valid := validation.Validation{}
	for k, v := range funcs {
		if err = validation.AddCustomFunc(k, v); err != nil {
			return http.StatusInternalServerError, e.UnknownError
		}
	}

	//表单验证
	check, err := valid.Valid(obj)
	if err != nil {
		logging.Error(err)
		return http.StatusInternalServerError, e.UnknownError
	}
	if !check {
		MakeErrors(valid.Errors)
		return http.StatusBadRequest, e.ParamInvalid
	}

	return http.StatusOK, e.OK
}
