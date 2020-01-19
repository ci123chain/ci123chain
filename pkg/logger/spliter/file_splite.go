package spliter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)


func NewFileLogger(fileDir, fileName string) FileLogger {
	//CloseLogger()
	f := FileLogger{
		fileDir: 		fileDir,
		fileName: 		fileName,
		mu:            	new(sync.RWMutex),
		stopTickerChan: make(chan bool, 1),
		errchan: 		make(chan string, 1),
	}

	t, _ := time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))
	f.date = &t
	go f.fileMonitor()
	return f
}

type FileLogger struct {
	fileDir        string         // 日志文件保存的目录
	fileName       string         // 日志文件名（无需包含日期和扩展名）
	logFile        *os.File       // 日志文件
	date           *time.Time     // 日志当前日期
	mu             *sync.RWMutex  // 读写锁，在进行日志分割和日志写入时需要锁住
	stopTickerChan chan bool      // 停止定时器的通道
	errchan 	   chan string
}

const DATE_FORMAT = "2006-01-02"


func (f *FileLogger) fileMonitor() {
	defer func() { recover() }()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if f.isMustSplit() {
				if err := f.split(); err != nil {
					if len(f.errchan) < 1 {
						f.errchan <- fmt.Sprintf("Log split error: %v", err)
					}
				}
			}
		case <-f.stopTickerChan:
			return
		}
	}
}

func (f *FileLogger)ErrChan() <- chan string{
	return f.errchan
}

func (f *FileLogger) isMustSplit() bool {
	t, _ := time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))
	return t.After(*f.date)
}

func (f *FileLogger) split() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	logFile := filepath.Join(f.fileDir, f.fileName)
	logFileBak := logFile + "-" + f.date.Format(DATE_FORMAT) + ".log"

	if f.logFile != nil {
		f.logFile.Close()
	}

	source, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(logFileBak)
	if err != nil {
		return err
	}
	_, err = io.Copy(destination, source)

	if err != nil {
		return err
	}

	t, _ := time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))

	f.date = &t

	os.Truncate(logFile, 0)


	//tempf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
	//if err != nil {
	//	return err
	//}
	//_, err = tempf.WriteAt([]byte(""), 0)
	return nil
}

