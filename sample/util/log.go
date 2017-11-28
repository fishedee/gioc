package util

import (
	"fmt"
	"github.com/fishedee/gioc"
)

var MyLog Log

type Log struct {
	Debug func(in string, args ...interface{})
	Info  func(in string, args ...interface{})
}

type logImpl struct {
}

func (this *logImpl) Debug(in string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+in+"\n", args...)
}

func (this *logImpl) Info(in string, args ...interface{}) {
	fmt.Printf("[INFO] "+in+"\n", args...)
}

func newLog() *logImpl {
	return &logImpl{}
}

func init() {
	MyLogImpl := newLog()
	MyLog.Debug = MyLogImpl.Debug
	MyLog.Info = MyLogImpl.Info
	gioc.Register(&Log{}, func() *logImpl {
		return MyLogImpl
	})
}
