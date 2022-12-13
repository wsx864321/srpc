package logger

import (
	"fmt"
	"testing"
)

func Test_getLineInfo(t *testing.T) {
	l := NewSweetLog()
	fmt.Println(l.caller())
}
