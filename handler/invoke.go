package handler

import (
	"ashe/protocol"
	"context"
	"fmt"
	"google.golang.org/grpc/peer"
)

type RpcService struct {
}

func (r *RpcService) Invoke(ctx context.Context, args *protocol.Args) (*protocol.Response, error) {
	var ip string
	if fromContext, ok := peer.FromContext(ctx); ok {
		network := fromContext.Addr.Network()
		ss := fromContext.Addr.String()
		fmt.Println(network, ss)
	}


	res := new(protocol.Response)
	client, err := protocol.NewGrpcClient(args.Version, args.Service, args.Method)
	if err != nil {
		res.Code = RemoteCallFailed
		res.Msg = InnerError.Error()
		return res, err
	}
	req := new(protocol.Request)
	req.RequestJson = args.Args
	fmt.Println("ip: ", ip)
	return client.Invoke(req, ip, nil)
}
