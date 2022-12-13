package logger

import (
	"context"
	"fmt"
	"log"
	"path"
	"runtime"
)

type Log interface {
	Debugf(ctx context.Context, format string, a ...interface{})
	Infof(ctx context.Context, format string, a ...interface{})
	Warnf(ctx context.Context, format string, a ...interface{})
	Errorf(ctx context.Context, format string, a ...interface{})
}

type SweetLog struct {
}

func NewSweetLog() *SweetLog {
	return &SweetLog{}
}

func (s *SweetLog) Debugf(ctx context.Context, format string, a ...interface{}) {
	log.Printf("【DEBUG】"+s.caller()+" "+format, a...)
}

func (s *SweetLog) Infof(ctx context.Context, format string, a ...interface{}) {
	log.Printf("【INFO】"+s.caller()+" "+format, a...)
}

func (s *SweetLog) Warnf(ctx context.Context, format string, a ...interface{}) {
	log.Printf("【WARN】"+s.caller()+" "+format, a...)
}

func (s *SweetLog) Errorf(ctx context.Context, format string, a ...interface{}) {
	log.Printf("【ERROR】"+s.caller()+" "+format, a...)
}

func (s *SweetLog) caller() string {
	var (
		pc       uintptr
		file     string
		lineNo   int
		ok       bool
		funcName string
	)

	pc, file, lineNo, ok = runtime.Caller(2)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}

	return fmt.Sprintf("%s/%s:%d", path.Base(file), path.Base(funcName), lineNo)
}
