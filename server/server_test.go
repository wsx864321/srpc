package server

import "testing"

func TestNewServer(t *testing.T) {
	s := NewServer()
	s.Start()
}
