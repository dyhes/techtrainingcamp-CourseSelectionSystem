package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogPath     = "../runtime/logs/" //日志文件的存放路径(相对于working directory的路径)
	LogFileName = "log"              //日志文件的名字
	LogFileExt  = "log"              //扩展名
	TimeFormat  = "20060102"         //日期格式
)

func getLogFilePath() string {
	return LogPath
}

//getLogFileFullPath 日志文件的路径包含文件名
func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogFileName, time.Now().Format(TimeFormat), LogFileExt)
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

//openLogFile 打开文件
func openLogFile(path string) *os.File {
	//得到文件信息结构描述文件
	_, err := os.Stat(path)
	switch {
	case os.IsNotExist(err):
		//如果日志文件不存在，那么就创建一个
		mkDir()
	case os.IsPermission(err):
		//没有权限
		log.Fatalf("Permission: %v", err)
	}
	//打开文件，文件权限为0644，即0rwxr--r--
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}
	return file
}

func mkDir() {
	//返回工作路径(wd stands for working directory)的根目录。
	//在这里，即desktop/course-selecting-sys
	dir, _ := os.Getwd()
	//MkdirAll是创建多级目录，Mkdir只能创建单级目录。os.ModePerm即目录权限为0777
	err := os.MkdirAll(dir+"/"+LogPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
