package logex

import (
	"net/http"
	"os"

	"context"
)

type Config struct {
	Dest  string `ha:"path"`
	Level string
}

func Agent(ctx context.Context, res http.ResponseWriter) context.Context {
	log := New()
	// TODO setloglevel
	log.SetLogLevel(DEBUG)
	// TODO setOutput
	log.SetOutput(os.Stdout)
	res.Header().Set("x-request-id", log.Id())
	ctx = context.WithValue(ctx, "logger", log)
	return ctx
}
