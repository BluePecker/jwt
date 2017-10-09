package httputils

import (
	"github.com/kataras/iris/context"
	"github.com/kataras/iris"
)

func Success(ctx context.Context, data interface{}) error {
	_, err := ctx.JSON(map[string]interface{}{
		"code":    iris.StatusOK,
		"data":    data,
		"message": "winner winner,chicken dinner.",
	})
	return err
}

func Failure(ctx context.Context, message string) error {
	_, err := ctx.JSON(map[string]interface{}{
		"code":    iris.StatusBadRequest,
		"data":    map[string]interface{}{},
		"message": message,
	})
	return err
}
