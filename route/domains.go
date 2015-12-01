package route

import (
	"config"
	"logex"
	"net/http"

	"golang.org/x/net/context"
)

type ds struct {
	Domains map[string]uint64 `yaml:"domains"`
}

func Domains(ctx context.Context, log logex.Logger, req *http.Request) (context.Context, error) {
	log.Debug("this is domain process")
	d := ds{}

	err := config.Load("routes.yaml", &d)
	if err != nil {
		return ctx, err
	}
	v, ok := d.Domains[req.Host]
	if !ok {
		v = 0
	}

	ctx = context.WithValue(ctx, "appid", v)
	return ctx, nil
}
