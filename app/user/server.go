package main

import (
	"ashe/app/user/api"
	"ashe/app/user/models"
	pb "ashe/app/user/proto"
	"ashe/common"
	"ashe/db"
	"flag"
	"fmt"
	"github.com/wuranxu/library/conf"
	nt "github.com/wuranxu/library/net"
	"github.com/wuranxu/library/service/etcd"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	config   = flag.String("configPath", `/Users/wuranxu/Downloads/ashe/config.json`, "config file path")
	yamlPath = flag.String("yamlPath", `/Users/wuranxu/Downloads/ashe/app/user/service.yaml`, "yaml file path")
	//config = flag.String("configPath", `G:\golang\ashe\config.json`, "config file path")
)

func main() {
	flag.Parse()
	common.Init(*config)
	var yamlConfig conf.YamlConfig
	if err := conf.ParseYaml(*yamlPath, &yamlConfig); err != nil {
		//if err := conf.ParseYaml(`./service.yaml`, &yamlConfig); err != nil {
		log.Fatal(err)
	}
	models.Conn = db.Init(models.Tables)
	port := fmt.Sprintf(":%d", yamlConfig.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("服务挂壁了, error: ", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterUserServer(s, &api.UserApi{})
	cli, err := etcd.NewClient(conf.Conf.Etcd)
	if err != nil {
		panic(err)
	}
	// 获取ip地址
	address := nt.GetLocalIp() + port
	err = cli.RegisterService(yamlConfig.Service, address, 5)
	if err != nil {
		panic(err)
	}
	if err := cli.RegisterApi(yamlConfig.Service, &api.UserApi{}, yamlConfig); err != nil {
		panic(err)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		cli.UnRegister(yamlConfig.Service, address)

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务意外关闭: %v", err)
	}
}
