package luxury

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/hjin-me/luxury/logex"

	"context"
)

func TestWrapEmpty(t *testing.T) {
	ctx := context.Background()
	fn := wrap(func() {
		t.Log("empty func")
	})
	var (
		tc context.Context
		ts string
		te error
	)
	tc, ts, te = fn(ctx)
	t.Log(tc, ts, te)
	if te != nil {
		t.Error("error not nil")
	}
	if ts != "default" {
		t.Error("step not default")
	}
}

func TestWrapFormat(t *testing.T) {
	ctx := context.Background()
	fn := wrap(func(ctx context.Context) (context.Context, string, error) {
		t.Log("normal fn")
		return ctx, "normal", errors.New("normal")
	})
	var (
		tc context.Context
		ts string
		te error
	)
	tc, ts, te = fn(ctx)
	t.Log(tc, ts, te)
	if te.Error() != "normal" {
		t.Error("error not normal")
	}
	if ts != "normal" {
		t.Error("step not normal")
	}
}
func TestLogger(t *testing.T) {
	log := logex.New()
	log.SetLogLevel(logex.DEBUG)
	bf := bytes.Buffer{}
	log.SetOutput(&bf)
	ctx := context.WithValue(context.Background(), "logger", log)
	fn := wrap(func(logger logex.Logger, ctx context.Context) context.Context {
		logger.Debug("123")
		t.Log("print log")
		t.Log(bf.String())
		return ctx
	})
	var (
		tc context.Context
		ts string
		te error
	)
	tc, ts, te = fn(ctx)
	t.Log(tc, ts, te)
	if strings.Index(bf.String(), "123") == -1 {
		t.Error(bf.String())
	}
}
