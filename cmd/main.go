package main

import (
	"fmt"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/setting"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/routers"
	"net/http"
)

func main() {
	router := routers.RegisterRouter()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort), //监听的端口号，格式为:8000
		Handler:        router,                               //http句柄，用于处理程序响应HTTP请求
		ReadTimeout:    setting.ReadTimeout,                  //读取超时时间
		WriteTimeout:   setting.WriteTimeout,                 //写超时时间
		MaxHeaderBytes: 1 << 20,                              //http报文head的最大字节数
	}

	server.ListenAndServe()
}
