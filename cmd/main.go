//package main
//
//import (
//	"ashe/common"
//	"ashe/db"
//	"ashe/models"
//	"ashe/router"
//	"flag"
//	"fmt"
//	"github.com/iris-contrib/middleware/cors"
//	_ "github.com/jinzhu/gorm"
//	"github.com/kataras/iris"
//	"github.com/kataras/iris/middleware/logger"
//	"github.com/kataras/iris/middleware/recover"
//)
//
//var (
//	configPath = flag.String("config", `./config.json`, "your config file with json")
//)
//
//func main() {
//	// parse arguments
//	flag.Parse()
//	// init config
//	common.Init(*configPath)
//	// init redis
//	common.InitRedisConnection()
//	// init database
//	db.Init()
//	// close db conn
//	defer models.Conn.Close()
//
//	app := iris.New()
//	app.AllowMethods(iris.MethodOptions, iris.MethodDelete) // cors
//	app.Use(recover.New())
//	app.Use(logger.New())
//	app.Use(cors.AllowAll())
//	router.AddHandler(app)
//	d, x, y := models.GetJobList(1, 20)
//	fmt.Println(len(d))
//	fmt.Println(models.DelJob(3))
//	//fmt.Println(models.NewAsheJob("吴冉旭爱洗澡332", "echo lixiaoyaoshiwo", "10.222.11.21", "wuranxu", 22))
//	d, x, y = models.GetJobList(1, 20)
//	fmt.Println(len(d))
//	fmt.Println(models.DelJob(3))
//	for _, h := range d {
//		fmt.Println(h)
//	}
//	app.Run(iris.Addr(":8088"), iris.WithoutServerError(iris.ErrServerClosed))
//}

package main

import (
	"ashe/proto/cronjob"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
}

func (s *Server) Add(c context.Context, j *cronjob.Job) (*cronjob.Response, error) {
	return &cronjob.Response{
		Code: 0,
		Msg:  "你好，应该成功了哦" + j.Name,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	cronjob.RegisterCronJobServiceServer(grpcServer, &Server{})
	grpcServer.Serve(lis)
}
