package logger

import (
	"log"
	"os"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

var AppLogger *Logger

func Init() *Logger {
	AppLogger = &Logger{
		infoLog:  log.New(os.Stdout, Green+"INFO: "+Reset, log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, Red+"ERROR: "+Reset, log.Ldate|log.Ltime|log.Lshortfile),
	}
	return AppLogger
}

func (l *Logger) Info(v ...interface{}) {
	l.infoLog.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.errorLog.Println(append([]interface{}{Red}, append(v, Reset)...)...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorLog.Printf(Red+format+Reset, v...)
}