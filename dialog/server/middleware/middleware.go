package middleware

import "github.com/kataras/iris/context"

var Provider []context.Handler = []context.Handler{}

func Register(handler context.Handler) {
	Provider = append(Provider, handler)
}
