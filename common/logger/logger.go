package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
	"log"
	"io"
	"os"
)

type logger struct {
	filename string
	fd       *lumberjack.Logger
}

var (
	ymLogger *logger
	lock     = new(sync.Mutex)
	once     = new(sync.Once)

	Info    = log.New(os.Stdout, "[Info] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[Warning] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error   = log.New(os.Stdout, "[Error] ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal   = log.New(os.Stdout, "[Fatal] ", log.Ldate|log.Ltime|log.Lshortfile)
)

func New(filename string, debug bool) *logger {
	if ymLogger != nil {
		return ymLogger
	}

	once.Do(func() {
		ymLogger = &logger{filename: filename}
		ymLogger.fd = &lumberjack.Logger{
			Filename: filename,
			MaxSize:  10,
			MaxAge:   28,
		}
		var w io.Writer
		if debug {
			w = io.MultiWriter(ymLogger.fd, os.Stdout)
		} else {
			w = ymLogger.fd
		}

		Info.SetOutput(w)
		Warning.SetOutput(w)
		Error.SetOutput(w)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	})
	return ymLogger
}

func Close() {
	if ymLogger == nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	if ymLogger.fd != nil {
		ymLogger.fd.Close()
	}
	ymLogger = nil
}

// TODO 封装?
//func Info(v ...interface{}) {
//	log.Println(" [info] ", fmt.Sprint(v...))
//}
//
//func Debug(v ...interface{}) {
//	log.Println(" [Debug] ", fmt.Sprint(v...))
//}
//
//func Error(v ...interface{}) {
//	log.Println(" [Error] ", fmt.Sprint(v...))
//}
///**
//	打印错误,并且
// */
//func Fatal(v ...interface{}) {
//	log.Println(" [Error] ", fmt.Sprint(v...))
//	os.Exit(1)
//}
