package main

import (
	ep "ashe/app/example"
	pb "ashe/app/example/proto"
	"ashe/common"
	"ashe/library/cache/etcd"
	"ashe/protocol"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"time"
)

func main() {

	common.Init("./config.json")
	client, err := etcd.NewClient(common.Conf.Etcd)
	if err != nil {
		panic(err)
	}
	r := etcd.NewResolver(client)
	resolver.Register(r)

	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", r.Scheme(), ep.ServiceName),
		grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cl := pb.NewHelloServiceClient(conn)

	for {
		resp, err := cl.Hello(context.Background(), &protocol.Request{RequestJson: "李逍遥"})
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			fmt.Println("res: ", resp.Msg)
		}
		<-time.After(time.Second)
	}
}
