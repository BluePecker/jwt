package middleware

import (
	"time"
	"github.com/Sirupsen/logrus"
	"strconv"
	"github.com/kataras/iris/context"
)

func init() {
	Register(func(ctx context.Context) {
		start := time.Now()
		ctx.Next()
		logrus.Infof("%v %4v %s %s %s", strconv.Itoa(ctx.GetStatusCode()), time.Now().Sub(start), ctx.RemoteAddr(), ctx.Method(), ctx.Path())
	})
}
