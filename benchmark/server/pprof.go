package main

import (
	"net/http"
	_ "net/http/pprof"
)

func pprof() {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", http.DefaultServeMux)
	}()
}
