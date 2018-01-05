package xlogger

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
)

// LogLevel is the log level type.
type LogLevel int

const (
	// DEBUG represents debug log level.
	DEBUG LogLevel = iota
	// INFO represents info log level.
	INFO
	// WARN represents warn log level.
	WARN
	// ERROR represents error log level.
	ERROR
	// FATAL represents fatal log level.
	FATAL
)

var tagName = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

type Logger struct {
	filename string
	fd       *lumberjack.Logger
	logfd    *log.Logger
	stdout   bool
	level    LogLevel
}

var (
	ymLogger *Logger
	lock     = new(sync.Mutex)
	once     = new(sync.Once)
)

func New(filename string, level LogLevel, stdout bool) *Logger {
	if ymLogger != nil {
		return ymLogger
	}

	once.Do(func() {
		ymLogger = &Logger{
			filename: filename,
			level:    level,
			stdout:   stdout,
		}
		ymLogger.fd = &lumberjack.Logger{
			Filename: filename,
			MaxSize:  10,
			MaxAge:   28,
		}
		var w io.Writer
		if stdout {
			w = io.MultiWriter(ymLogger.fd, os.Stdout)
		} else {
			w = ymLogger.fd
		}
		ymLogger.logfd = log.New(os.Stdout, "", log.LstdFlags)
		ymLogger.logfd.SetOutput(w)
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
func getRuntimeInfo() (string, string, int) {
	pc, fn, ln, ok := runtime.Caller(3) // 3 steps up the stack frame
	if !ok {
		fn = "???"
		ln = 0
	}
	function := "???"
	caller := runtime.FuncForPC(pc)
	if caller != nil {
		function = caller.Name()
	}
	return function, fn, ln
}

func (l Logger) doPrintln(level LogLevel, v ...interface{}) {
	if l.logfd == nil {
		return
	}

	if level >= l.level {
		funcName, fileName, lineNum := getRuntimeInfo()
		prefix := fmt.Sprintf("[%5s] [%s] (%s:%d) - ", tagName[level], path.Base(funcName),
			path.Base(fileName), lineNum)
		value := fmt.Sprintf("%s %s", prefix, fmt.Sprintln(v...))
		l.logfd.Print(value)

		if level == FATAL {
			os.Exit(1)
		}
	}

}

func (l Logger) doPrintf(level LogLevel, format string, v ...interface{}) {
	if l.logfd == nil {
		return
	}
	if level >= l.level {
		funcName, fileName, lineNum := getRuntimeInfo()
		format = fmt.Sprintf("%5s [%s] (%s:%d) - %s", tagName[level], path.Base(funcName), path.Base(fileName), lineNum, format)
		l.logfd.Printf(format, v...)
		if level == FATAL {
			os.Exit(1)
		}
	}

}

func Infof(format string, v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志未初始化: ", v)
		return
	}
	ymLogger.doPrintf(INFO, format, v...)
}
func Info(v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志未初始化: ", v)
		return
	}
	ymLogger.doPrintln(INFO, v...)

}
func Debugf(format string, v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志未初始化: ", v)
		return
	}
	ymLogger.doPrintf(DEBUG, format, v...)
}

func Debug(v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志未初始化: ", v)
		return
	}
	ymLogger.doPrintln(DEBUG, v...)
}

func Warnf(format string, v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志未初始化: ", v)
		return
	}
	ymLogger.doPrintf(WARN, format, v...)
}
func Warn(v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志未初始化: ", v)
		return
	}
	ymLogger.doPrintln(WARN, v...)
}

func Errorf(format string, v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志未初始化: ", v)
		return
	}
	ymLogger.doPrintf(ERROR, format, v...)
}

func Error(v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志未初始化: ", v)
		return
	}
	ymLogger.doPrintln(ERROR, v...)
}

/**
打印错误,并且
*/
func Fatal(v ...interface{}) {
	if ymLogger == nil {
		fmt.Println("日志为初始化: ", v)
		return
	}
	ymLogger.doPrintln(FATAL, v...)
}
