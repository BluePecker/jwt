package router

import "github.com/kataras/iris"

type (
	Router interface {
		Routes(engine *iris.Application)
	}
)
