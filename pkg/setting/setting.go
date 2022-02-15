package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"time"
)

//var (
//	Config *ini.File
//
//	RunMode string
//
//	PageSize  int
//	JWTSecret string
//
//	HttpPort     int
//	ReadTimeout  time.Duration
//	WriteTimeout time.Duration
//
//	User        string
//	Password    string
//	Host        string
//	Name        string
//	TablePrefix string
//)

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

////init 包导入的时候执行以读取ini文件
//func init() {
//	var err error
//	Config, err = ini.Load("../conf/config.ini")
//	if err != nil {
//		log.Fatalln(fmt.Sprintf("Read ini file failed: %v", err))
//	}
//	loadBase()
//	loadApp()
//	loadServer()
//	loadDatabase()
//}
//
//func loadBase() {
//	RunMode = Config.Section("").Key("RUN_MODE").MustString("debug")
//}
//
//func loadApp() {
//	app, err := Config.GetSection("app")
//	if err != nil {
//		logging.Fatal(fmt.Sprintf("Can't read section [app]: %v", err))
//	}
//	JWTSecret = app.Key("JWT_SECRET").String()
//	PageSize = app.Key("PAGE_SIZE").MustInt(10)
//}
//
////loadServer 读取ini文件中的[server] section
//func loadServer() {
//	server, err := Config.GetSection("server")
//	if err != nil {
//		logging.Fatal(fmt.Sprintf("Can't read section [server]: %v", err))
//	}
//	HttpPort = server.Key("HTTP_PORT").MustInt(8000)
//	ReadTimeout = time.Duration(server.Key("READ_TIMEOUT").MustInt(60)) * time.Second
//	WriteTimeout = time.Duration(server.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
//}
//
////loadDatabase 读取ini文件中的[database] section
//func loadDatabase() {
//	db, err := Config.GetSection("database")
//	if err != nil {
//		logging.Fatal(fmt.Sprintf("Can't read section [database]: %v", err))
//	}
//	User = db.Key("USER").String()
//	Password = db.Key("PASSWORD").String()
//	Host = db.Key("HOST").String()
//	Name = db.Key("NAME").String()
//	TablePrefix = db.Key("TABLE_PREFIX").String()
//}
