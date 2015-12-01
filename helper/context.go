package helper

import "golang.org/x/net/context"

func AppId(ctx context.Context) (uint64, bool) {
	v, ok := ctx.Value("appid").(uint64)
	return v, ok
}
