package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

//Level 日志级别
type Level int

const (
	INFO    = iota //信息
	DEBUG          //调试
	WARNING        //警告
	ERROR          //错误
	FATAL          //失败
)

var (
	File *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{INFO: "INFO", DEBUG: "DEBUG", WARNING: "WARNING", ERROR: "ERROR", FATAL: "FATAL"}
)

func init() {
	filePath := getLogFileFullPath()
	File = openLogFile(filePath)
	//func New(out io.Writer, prefix string, flag int) *Logger函数创建一个新的日志记录器。
	//out为要写入日志数据的IO句柄；prefix定义每个生成的日志行的开头；flag定义了日志记录的属性。
	//其他属性有：
	//const (
	//    Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	//    Ltime                         // the time in the local time zone: 01:23:23
	//    Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	//    Llongfile                     // full file name and line number: /a/b/c/d.go:23
	//    Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	//    LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	//    LstdFlags     = Ldate | Ltime // initial values for the standard logger
	//)
	logger = log.New(File, DefaultPrefix, log.LstdFlags)
}

//setPrefix 设置日志行的前缀，主要是日志级别，输出日志的函数所在文件以及行数等信息
func setPrefix(level Level) {
	//func Caller(skip int) (pc uintptr, file string, line int, ok bool)
	//skip是要提升的堆栈帧数。0表示当前函数，1表示上一层函数，2表示上上层函数，依次递增
	//pc是函数指针; file是函数所在文件名; line为所在行号; ok表示是否可以获取到信息
	//DefaultCallerDepth的值为2
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		//filepath.Base返回文件路径名的最后一个元素。例如/foo/bar/baz.js返回baz.js
		//日志前缀的默认格式为[级别][函数所在的文件:函数所在的行号]
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Warning(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}
