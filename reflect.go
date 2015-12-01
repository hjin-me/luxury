package luxury

import (
	"logex"
	"net/http"
	"reflect"

	"golang.org/x/net/context"
)

type formatAgent func(context.Context) (context.Context, string, error)

type Agent interface{}

var (
	iContext  []context.Context
	iError    []error
	iLogger   []logex.Logger
	iRequest  []*http.Request
	iResponse []http.ResponseWriter
)

func wrap(fn Agent) formatAgent {
	r := reflect.ValueOf(fn)
	t := r.Type()
	in := make([]reflect.Value, t.NumIn())
	return func(ctx context.Context) (context.Context, string, error) {
		for i := 0; i < t.NumIn(); i++ {
			switch {
			case t.In(i).Implements(reflect.TypeOf(iContext).Elem()):
				in[i] = reflect.ValueOf(ctx)
			case t.In(i).Implements(reflect.TypeOf(iLogger).Elem()):
				logger, ok := ctx.Value("logger").(logex.Logger)
				if !ok {
					logger = logex.New()
				}
				in[i] = reflect.ValueOf(logger)
			case t.In(i).AssignableTo(reflect.TypeOf(iRequest).Elem()):
				v, ok := ctx.Value("req").(*http.Request)
				if !ok {
					panic("*http.Request not found")
				}
				in[i] = reflect.ValueOf(v)
			case t.In(i).Implements(reflect.TypeOf(iResponse).Elem()):
				v, ok := ctx.Value("res").(http.ResponseWriter)
				if !ok {
					panic("http.ResponseWriter not found")
				}
				in[i] = reflect.ValueOf(v)
			}
		}
		/*
			if t.NumIn() == 1 && reflect.TypeOf(ctx).AssignableTo(t.In(0)) {
				in[0] = reflect.ValueOf(ctx)
			}
		*/
		out := r.Call(in)
		var (
			nextStep string = "default"
			err      error
		)
		for _, v := range out {
			switch {
			case v.Type().Implements(reflect.TypeOf(iContext).Elem()):
				// log.Println("got context.Context")
				ctx = v.Interface().(context.Context)
			case v.Type().Implements(reflect.TypeOf(iError).Elem()) && !v.IsNil():
				// log.Println("got error")
				err = v.Interface().(error)
			case v.Type().AssignableTo(reflect.TypeOf(nextStep)):
				// log.Println("got step")
				nextStep = v.String()
			default:
				// log.Println(v.Type())
			}
		}
		return ctx, nextStep, err
	}
}
