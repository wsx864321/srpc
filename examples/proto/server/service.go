package main

import (
	"context"
	"fmt"
)

type HelloWorld struct {
}

func (h *HelloWorld) SayHello(ctx context.Context, req *HelloWorldReq) (*HelloWorldResp, error) {
	fmt.Println(req.Name)
	return &HelloWorldResp{
		Msg: fmt.Sprintf("%s say hello", req.Name),
	}, nil
}
