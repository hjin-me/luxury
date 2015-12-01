package route

import (
	"helper"
	"log"

	"golang.org/x/net/context"
)

func Url(ctx context.Context) (context.Context, string, error) {
	id, ok := helper.AppId(ctx)
	if !ok {
		return ctx, "404", nil
	}
	log.Println(id)
	switch id {
	case 1:
		ctx = context.WithValue(ctx, "theme", "cp")
		return ctx, "cp_starter", nil
	case 2:
		ctx = context.WithValue(ctx, "theme", "my")
	case 3:
		ctx = context.WithValue(ctx, "theme", "content")
	default:
		return ctx, "404", nil
	}
	_ = id
	return ctx, "post", nil
}
