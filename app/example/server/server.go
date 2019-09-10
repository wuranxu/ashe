package main

import (
	"ashe/app/example"
	pb "ashe/app/example/proto"
	"ashe/common"
	"ashe/library/cache/etcd"
	nt "ashe/library/net"
	"ashe/protocol"
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	port = flag.String("port", example.Port, "grpc endpoints")
	conf = flag.String("configPath", "../../../config.json", "config file path")
)

type hello struct{}

func (h *hello) Hello(ctx context.Context, request *protocol.Request) (*protocol.Response, error) {
	return &protocol.Response{Msg: "你好, " + request.RequestJson}, nil
}

func (h *hello) Assert(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	result := new(protocol.Response)
	if in.RequestJson == "334" {
		result.Code = 200
		return result, nil
	}
	result.Code = 5000
	result.Msg = "并不相等哦"
	result.ResultJson = `{"version": "男丁啊"}`
	return result, nil
}



func main() {
	flag.Parse()
	common.Init(*conf)
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal("服务挂壁了, error: ", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &hello{})
	cli, err := etcd.NewClient(common.Conf.Etcd)
	if err != nil {
		panic(err)
	}
	err = cli.RegisterService(example.ServiceName, nt.GetLocalIp()+*port, 1)
	if err != nil {
		panic(err)
	}
	if err := cli.RegisterApi(example.ServiceName, &hello{}); err != nil {
		panic(err)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		cli.UnRegister(example.ServiceName, nt.GetLocalIp()+*port)

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("fff")
	}
}
