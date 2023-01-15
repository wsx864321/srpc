package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wsx864321/srpc/pool"
)

const addr string = "127.0.0.1:7777"

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	go server()
	//等待tcp server启动
	time.Sleep(2 * time.Second)
	client()
	fmt.Println("使用: ctrl+c 退出服务")
	<-c
	fmt.Println("服务退出")
}

func client() {

	p, err := pool.NewChannelPool()
	if err != nil {
		fmt.Println("err=", err)
	}

	//从连接池中取得一个连接
	v, err := p.Get(context.Background(), "tcp", addr)

	//将连接放回连接池中
	p.Put(v)

	//释放连接池中的所有连接
	//p.Release()

	//查看当前连接中的数量
	current := p.Len()
	fmt.Println("len=", current)
}

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on ", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		//go handleRequest(conn)
	}
}
