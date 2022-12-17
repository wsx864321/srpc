package main

import (
	"context"

	"github.com/wsx864321/srpc/metadata"
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
	metaData := metadata.ExtractServerMetadata(ctx)
	return &HelloWorldResp{
		Msg: metaData["token"],
	}, nil
}
