package util

import (
	"fmt"
	"github.com/fishedee/gioc"
)

var MyLog Log

type Log interface {
	Debug(in string, args ...interface{})
	Info(in string, args ...interface{})
}

type logImpl struct {
}

func (this *logImpl) Debug(in string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+in+"\n", args...)
}

func (this *logImpl) Info(in string, args ...interface{}) {
	fmt.Printf("[INFO] "+in+"\n", args...)
}

func newLog() Log {
	return &logImpl{}
}

func init() {
	MyLog = newLog()
	gioc.Register(func() Log {
		return MyLog
	})
}
