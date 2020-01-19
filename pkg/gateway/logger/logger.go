package logger

import (
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/logger/spliter"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)


type LEVEL byte

const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
)

type FileLogger struct {
	fileDir        string         // 日志文件保存的目录
	fileName       string         // 日志文件名（无需包含日期和扩展名）

	logLevel       LEVEL          // 日志等级
	logFile        *os.File       // 日志文件
	logChan        chan string    // 日志消息通道，以实现异步写日志

	lg             *log.Logger    // 系统日志对象
}


var fileLogger *FileLogger
var spliteLogger spliter.FileLogger

func Init(fileDir, fileName, prefix, level string) error {
	CloseLogger()

	f := &FileLogger{
		fileDir:       os.ExpandEnv(fileDir),
		fileName:      fileName,
		logChan:       make(chan string, 1),
	}

	switch strings.ToUpper(level) {
	case "DEBUG":
		f.logLevel = DEBUG
	case "WARN":
		f.logLevel = WARN
	case "ERROR":
		f.logLevel = ERROR
	default:
		f.logLevel = INFO
	}


	if f.fileName == "" {
		f.lg = log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds)
	} else {
		file, err := f.CreateFile()
		if err != nil {
			return err
		}
		f.logFile = file

		f.lg = log.New(io.MultiWriter(f.logFile, os.Stdout), prefix, log.LstdFlags|log.Lmicroseconds)
	}

	spliteLogger = spliter.NewFileLogger(f.fileDir, fileName)
	fileLogger = f


	go f.logWriter()
	go errWriter()

	return nil
}

func CloseLogger() {
	if fileLogger != nil {
		close(fileLogger.logChan)
		fileLogger.lg = nil
		fileLogger.logFile.Close()
	}
}

func (f *FileLogger) CreateFile() (*os.File, error) {
	if _, err := os.Stat(f.fileDir); os.IsNotExist(err) {
		os.MkdirAll(f.fileDir, os.ModePerm)
		os.Chmod(f.fileDir, os.ModePerm)
	}
	fullFileName := filepath.Join(f.fileDir, f.fileName)
	file, err := os.OpenFile(fullFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	return file, err
}

func errWriter()  {
	defer func() { recover() }()
	for {
		str, ok := <- spliteLogger.ErrChan()
		if !ok {
			return
		}
		Error(str)
	}
}

func (f *FileLogger) logWriter() {
	defer func() { recover() }()

	for {
		str, ok := <-f.logChan
		if !ok {
			return
		}

		f.lg.Output(2, str)
	}
}


func Error(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLogger.logLevel <= ERROR {
		fileLogger.logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[Error]"+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLogger.logLevel <= WARN {
		fileLogger.logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[Warning]"+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLogger.logLevel <= INFO {
		fileLogger.logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[INFO]"+format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLogger.logLevel <= DEBUG {
		fileLogger.logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[Debug]"+format, v...)
	}
}

