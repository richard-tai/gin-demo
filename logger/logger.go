package logger

// 支持文件输出
// 支持控制台输出
// 支持等级控制
// 支持日志文件滚动
// 不支持不同等级分文件输出

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
)

var D *Logger = New(Conf{enableConsole: true})

func New(cf Conf) *Logger {
	if cf.flags == 0 {
		cf.flags = DEFAULT_FLAGS
	}
	lg := &Logger{
		conf:        cf,
		curFilePath: path.Join(cf.dirPath, cf.fileName),
	}
	if cf.enableFile {
		if err := os.MkdirAll(lg.curFilePath, 0777); err != nil {
			log.Printf("%v", err)
		}
		lg.createFileLogger()
	}
	if cf.enableConsole {
		log.SetFlags(cf.flags)
	}
	return lg
}

type Conf struct {
	dirPath       string
	fileName      string
	fileNumLimit  int
	fileByteLimit int64
	prefix        string
	flags         int
	level         LEVEL
	enableFile    bool
	enableConsole bool
}

type Logger struct {
	conf          Conf
	curFilePath   string
	curFile       *os.File
	curFileLogger *log.Logger
	sync.RWMutex
}

func (l *Logger) Trace(format string, args ...interface{}) {
	l.out(TRACE, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.out(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.out(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.out(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.out(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.out(FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) out(lv LEVEL, format string, args ...interface{}) {
	defer catchError()
	if lv >= l.conf.level {
		outStr := l.conf.prefix + " " + mapLevelTag[lv] + " " + format
		if len(args) > 0 {
			outStr = fmt.Sprintf(outStr, args...)
		}
		if l.conf.enableFile {
			l.checkLogFile()
			l.RLock()
			defer l.RUnlock()
			if l.curFileLogger != nil {
				l.curFileLogger.Output(3, outStr)
			} else {
				log.Printf("file logger is nil")
			}
		}
		if l.conf.enableConsole {
			log.Output(3, outStr)
		}
	}
}

func (l *Logger) checkLogFile() {
	if l.conf.fileNumLimit > 1 && l.conf.fileByteLimit > 0 {
		fz := fileSize(l.curFilePath)
		if fz >= l.conf.fileByteLimit || fz < -1 {
			l.Lock()
			defer l.Unlock()
			if fz >= l.conf.fileByteLimit {
				l.rotateFile()
			}
			l.createFileLogger()
		}
	}
}

// xx -> xx.1 -> xx.2
func (l *Logger) rotateFile() {
	lastFilePath := l.curFilePath + strconv.Itoa(l.conf.fileNumLimit-1)
	if err := os.Remove(lastFilePath); err != nil {
		log.Printf("%v", err)
	}
	for i := l.conf.fileNumLimit - 2; i >= 1; i-- {
		newPath := l.curFilePath + strconv.Itoa(i+1)
		oldPath := l.curFilePath + strconv.Itoa(i)
		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("%v", err)
		}
	}
	secondFilePath := l.curFilePath + strconv.Itoa(1)
	if err := os.Rename(l.curFilePath, secondFilePath); err != nil {
		log.Printf("%v", err)
	}
}

func fileSize(fp string) int64 {
	info, err := os.Stat(fp)
	if err != nil {
		log.Printf("%v", err)
		return -1
	}
	return info.Size()
}

func (l *Logger) createFileLogger() error {
	if l.curFile != nil {
		l.curFile.Close()
		l.curFile = nil
	}
	var err error
	l.curFile, err = os.OpenFile(l.curFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err == nil && l.curFile != nil {
		l.curFileLogger = log.New(l.curFile, "", l.conf.flags)
	} else {
		log.Printf("%v, %v", err, l.curFile)
	}
	return err
}

func catchError() {
	if err := recover(); err != nil {
		log.Printf("%v", err)
	}
}

type LEVEL int

const (
	ALL   LEVEL = iota // 0
	TRACE              // 1
	DEBUG              // 2
	INFO               // 3
	WARN               // 4
	ERROR              // 5
	FATAL              // 6
	OFF                // 7
)

var mapLevelTag = map[LEVEL]string{
	TRACE: "[trace]",
	DEBUG: "[debug]",
	INFO:  "[info]",
	WARN:  "[warn]",
	ERROR: "[error]",
	FATAL: "[fatal]",
}

var DEFAULT_FLAGS = log.Lshortfile | log.Ldate | log.Lmicroseconds
