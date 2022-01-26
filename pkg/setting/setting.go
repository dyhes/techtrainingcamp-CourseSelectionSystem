package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"time"
)

type App struct {
	RunMode string
	Host    string //服务器公网IP
}

type Server struct {
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	User     string
	Password string
	Host     string
	Name     string
}

var (
	config *ini.File

	AppSetting      = &App{}
	ServerSetting   = &Server{}
	DatabaseSetting = &Database{}
)

//读取配置文件并映射section到struct
func Setup() {
	//读取配置文件
	var err error
	config, err = ini.Load("../conf/config.ini")
	if err != nil {
		log.Fatalln(fmt.Sprintf("Read ini file failed: %v", err))
	}

	//映射section到struct
	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)

	//这里必须要乘一个time.Second否则默认是纳秒
	ServerSetting.ReadTimeout *= time.Second
	ServerSetting.WriteTimeout *= time.Second
}

func mapTo(section string, v interface{}) {
	err := config.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("config.MapTo %s err: %v", section, err)
	}
}
