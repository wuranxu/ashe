package main

import (
	"ashe/app/user/api"
	"ashe/app/user/models"
	pb "ashe/app/user/proto"
	"ashe/common"
	"ashe/db"
	"ashe/library/cache/etcd"
	"ashe/library/conf"
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
	ServiceName = "user"
	Port        = ":12788"
)

var (
	port    = flag.String("port", Port, "grpc endpoints")
	service = flag.String("service", ServiceName, "grpc endpoints")
	//config    = flag.String("configPath", `../../config.json`, "config file path")
	//config = flag.String("configPath", `/Users/wuranxu/Downloads/ashe/config.json`, "config file path")
	config = flag.String("configPath", `G:\golang\ashe\config.json`, "config file path")
)

func main() {
	flag.Parse()
	common.Init(*config)
	models.Conn = db.Init(models.Tables)
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal("服务挂壁了, error: ", err)
	}
	var yamlConfig common.YamlConfig
	if err := conf.ParseYaml(`G:\golang\ashe\app\user\service.yaml`, &yamlConfig); err != nil {
	//if err := conf.ParseYaml(`./service.yaml`, &yamlConfig); err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterUserServer(s, &api.UserApi{})
	cli, err := etcd.NewClient(common.Conf.Etcd)
	if err != nil {
		panic(err)
	}
	err = cli.RegisterService(*service, nt.GetLocalIp()+*port, 5)
	if err != nil {
		panic(err)
	}
	if err := cli.RegisterApi(*service, &api.UserApi{}, yamlConfig); err != nil {
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
