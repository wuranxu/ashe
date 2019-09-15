package main

import (
	"ashe/common"
	"ashe/handler"
	"ashe/library/auth"
	"github.com/kataras/iris"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	common.Init("./config.json")
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
	app.Handle("POST", "/api/:service/:method", auth.Auth(handler.CallRpcWithAuth))
	app.Handle("POST", "/auth/:service/:method", handler.CallRpc)

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
