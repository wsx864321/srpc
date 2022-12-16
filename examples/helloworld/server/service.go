package main

import (
	"context"
	"fmt"
)

type HelloWorld struct {
}

type HelloWorldReq struct {
	Name string `json:"name"`
}

type HelloWorldResp struct {
	Msg string `json:"msg"`
}

func (h *HelloWorld) SayHello(ctx context.Context, req *HelloWorldReq) (*HelloWorldResp, error) {
	return &HelloWorldResp{
		Msg: fmt.Sprintf("%s say hello", req.Name),
	}, nil
}
