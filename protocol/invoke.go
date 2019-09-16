package protocol

import (
	"ashe/common"
	"ashe/library/auth"
	"ashe/library/cache/etcd"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"time"
)

var (
	MethodNotFound = errors.New("没有找到对应的方法，请检查您的参数")
)

type GrpcClient struct {
	cc     *grpc.ClientConn
	cli    *etcd.Client
	method etcd.Method
}

func Call(service, method string, in *Request, opt ...grpc.CallOption) (*Response, error) {
	client, err := NewGrpcClient(service, method)
	if err != nil {
		return nil, err
	}
	return client.Invoke(in, nil, opt...)
}

func (c *GrpcClient) Invoke(in *Request, userInfo *auth.CustomClaims, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if userInfo != nil {
		md := metadata.New(map[string]string{"user": userInfo.Marshal()})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	if err := c.cc.Invoke(ctx, c.method.Path, in, out, opts...); err != nil {
		return out, err
	}
	return out, nil
}

func (c *GrpcClient) Auth() bool {
	return c.method.Auth
}

func (c *GrpcClient) Close() error {
	if c != nil {
		return c.cc.Close()
	}
	return nil
}

func (c *GrpcClient) getCallAddr(service, method string) (etcd.Method, error) {
	var md etcd.Method
	addr := c.cli.GetSingle(fmt.Sprintf("%s.%s", service, method))
	if addr == "" {
		return md, MethodNotFound
	}
	if err := json.Unmarshal([]byte(addr), &md); err != nil {
		return md, err
	}
	return md, nil
}

func NewGrpcClient(service, method string) (*GrpcClient, error) {
	cl, err := etcd.NewClient(common.Conf.Etcd)
	if err != nil {
		return nil, err
	}
	re := etcd.NewResolver(cl)
	resolver.Register(re)
	// 3秒未连接上直接返回
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", re.Scheme(), service),
		grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := &GrpcClient{cli: cl, cc: conn}
	if client.method, err = client.getCallAddr(service, method); err != nil {
		return nil, err
	}
	return client, nil
}
