package main

import (
	"ashe/common"
	"ashe/handler"
	"flag"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

var (
	serverHost = flag.String("host", "0.0.0.0", "网关服务地址")
	serverPort = flag.Int("port", 8080, "网关服务端口号")
	configPath = flag.String("config", "./config.json", "网关配置文件")
)

func main() {
	flag.Parse()
	common.Init(*configPath)
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Method:   GET
	// Resource: http://localhost:8080
	//app.Handle("POST", "/:service/:method", auth.Auth(handler.CallRpc))
	//app.Handle("POST", "/api/:service/:method", auth.Auth(handler.CallRpcWithAuth))
	app.Handle("POST", "/:service/:method", handler.CallRpc)

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})
	app.Run(iris.Addr(fmt.Sprintf("%s:%d", *serverHost, *serverPort)), iris.WithoutServerError(iris.ErrServerClosed))
}
