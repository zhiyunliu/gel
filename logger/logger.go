package logger

import (
	"sync"
)

type Logger interface {
	Debug(content ...interface{})
	Debugf(format string, content ...interface{})
	Info(content ...interface{})
	Infof(format string, content ...interface{})
	Error(content ...interface{})
	Errorf(format string, content ...interface{})
	Fatal(content ...interface{})
	Fatalf(format string, content ...interface{})
	Close()
}

type internalLogger interface {
	Logger
	SetSessionID(sessionId string)
}

var loggerPool *sync.Pool
var closeChan chan struct{}

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return New()
		},
	}
	closeChan = make(chan struct{})
}

func New() Logger {
	return newZeroLogger()
}

func GetLogger(sessionID string) (logger Logger) {
	tmpLogger := loggerPool.Get().(internalLogger)
	tmpLogger.SetSessionID(sessionID)
	return tmpLogger
}