package app

import (
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

//发送http响应报文
func (g *Gin) Response(httpCode int, errCode e.ErrNo, data interface{}) {
	if data == nil {
		g.C.JSON(httpCode, gin.H{
			"Code": errCode,
			"Msg":  e.GetMsg(errCode),
			//"data": data,
		})
	} else {
		g.C.JSON(httpCode, gin.H{
			"Code": errCode,
			"Msg":  e.GetMsg(errCode),
			"Data": data,
		})
	}
}
