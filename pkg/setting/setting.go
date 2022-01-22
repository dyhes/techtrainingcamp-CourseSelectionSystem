package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/hollowdjj/course-selecting-sys/pkg/logging"
	"log"
	"time"
)

var (
	Config *ini.File

	RunMode string

	JWTSecret string

	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
)

//init 包导入的时候执行以读取ini文件
func init() {
	var err error
	Config, err = ini.Load("./conf/config.ini")
	if err != nil {
		log.Fatalln(fmt.Sprintf("Read ini file failed: %v", err))
	}
	loadBase()
	loadApp()
	loadServer()
	loadDatabase()
}

func loadBase() {
	RunMode = Config.Section("").Key("RUN_MODE").MustString("debug")
}

func loadApp() {
	app, err := Config.GetSection("app")
	if err != nil {
		logging.Fatal(fmt.Sprintf("Can't read section [app]: %v", err))
	}
	JWTSecret = app.Key("JWT_SECRET").String()
}

//loadServer 读取ini文件中的[server] section
func loadServer() {
	server, err := Config.GetSection("server")
	if err != nil {
		logging.Fatal(fmt.Sprintf("Can't read section [server]: %v", err))
	}
	HttpPort = server.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(server.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(server.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

//loadDatabase 读取ini文件中的[database] section
func loadDatabase() {
	db, err := Config.GetSection("database")
	if err != nil {
		logging.Fatal(fmt.Sprintf("Can't read section [database]: %v", err))
	}
	User = db.Key("USER").String()
	Password = db.Key("PASSWORD").String()
	Host = db.Key("HOST").String()
	Name = db.Key("NAME").String()
	TablePrefix = db.Key("TABLE_PREFIX").String()
}
