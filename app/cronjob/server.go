package main

import (
	"ashe/app/cronjob/api"
	"ashe/app/cronjob/models"
	pb "ashe/app/cronjob/proto"
	"ashe/common"
	"ashe/db"
	"ashe/library/cache/etcd"
	"ashe/library/cronjob"
	nt "ashe/library/net"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	ServiceName = "cronjob"
	Port        = ":12688"
)

var (
	port    = flag.String("port", Port, "grpc endpoints")
	service = flag.String("service", ServiceName, "grpc endpoints")
	//conf    = flag.String("configPath", `../../config.json`, "config file path")
	conf = flag.String("configPath", `/Users/wuranxu/Downloads/ashe/config.json`, "config file path")
)

func main() {
	flag.Parse()
	common.Init(*conf)
	models.Conn = db.Init(models.Tables)
	cronjob.InitRedisConnection(common.Conf.Redis)
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal("服务挂壁了, error: ", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterCronjobServer(s, &api.Job{})
	cronjob.InitRedisConnection(common.Conf.Redis)
	cli, err := etcd.NewClient(common.Conf.Etcd)
	if err != nil {
		panic(err)
	}
	err = cli.RegisterService(*service, nt.GetLocalIp()+*port, 5)
	if err != nil {
		panic(err)
	}
	if err := cli.RegisterApi(*service, &api.Job{}); err != nil {
		panic(err)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		cli.UnRegister(*service, nt.GetLocalIp()+*port)

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
