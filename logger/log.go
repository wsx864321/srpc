package logger

import (
	"context"
	"log"
)

type Log interface {
	Debugf(ctx context.Context, format string, a ...interface{})
	Infof(ctx context.Context, format string, a ...interface{})
	Warnf(ctx context.Context, format string, a ...interface{})
	Errorf(ctx context.Context, format string, a ...interface{})
}

type SweetLog struct {

}

func NewSweetLog() *SweetLog  {
	return &SweetLog{}
}

func (s *SweetLog) Debugf(ctx context.Context, format string, a ...interface{}) {
	log.Printf("【DEBUG】" + format, a...)
}

func (s *SweetLog) Infof(ctx context.Context, format string, a ...interface{}) {
	log.Printf("【INFO】" + format, a...)
}

func (s *SweetLog) Warnf(ctx context.Context, format string, a ...interface{}) {
	log.Printf("【WARN】" + format, a...)
}

func (s *SweetLog) Errorf(ctx context.Context, format string, a ...interface{}) {
	log.Printf("【ERROR】" + format, a...)
}