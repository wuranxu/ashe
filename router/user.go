package router

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)


func AddHandler(app *iris.Application) {
	app.Handle("GET", "/", func(context context.Context) {
		context.Write([]byte(`{"message": "你好"}`))
	})
}