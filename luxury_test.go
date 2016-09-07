package luxury

import (
	"log"
	"path/filepath"
	"testing"

	"context"
)

func TestWorkflow(t *testing.T) {
	testConf, err := filepath.Abs("./test/workflow.yaml")
	if err != nil {
		t.Fatal(err)
	}

	counter := 0
	wf := New()

	ok := func(ctx context.Context) (context.Context, string, error) {
		counter = counter + 1
		return ctx, "OK", nil
	}
	NewAgent(wf.Context, ok, "start", "step1", "step2", "step3", "step4")

	err = wf.Load(testConf)
	if err != nil {
		t.Log(testConf)
		t.Error(err)
	}
	err = wf.Handle("start")
	if err != nil {
		t.Error(err)
	}
	if counter != 3 {
		t.Error(counter, "is not 3")
	}

	counter = 0
	err = wf.Handle("step1")
	if err != nil {
		t.Error(err)
	}
	if counter != 1 {
		t.Error(counter, "is not 1")
	}
}

func TestContextValueChange(t *testing.T) {
	testConf, err := filepath.Abs("./test/workflow.yaml")
	if err != nil {
		t.Fatal(err)
	}

	okFn := func(ctx context.Context) (context.Context, string, error) {
		ctx = context.WithValue(ctx, "key", "value")
		return ctx, "OK", nil
	}
	wf := New()
	NewAgent(wf.Context, okFn, "start", "step1", "step2", "step3", "step4")

	err = wf.Load(testConf)
	if err != nil {
		t.Log(testConf)
		t.Error(err)
	}
	err = wf.Handle("start")
	if err != nil {
		t.Error(err)
	}
	value, ok := wf.Context.Value("key").(string)
	t.Log(value, ok)
	if !ok {
		t.Error("value is not string")
	}
	if value != "value" {
		t.Error("value is not value", value)
	}

}

func TestDynamicStep(t *testing.T) {

	log.Println("test dynamic")
	testConf, err := filepath.Abs("./test/dynamic_step.yaml")
	if err != nil {
		t.Fatal(err)
	}

	starterFn := func() {
	}

	dynamicFn := func(ctx context.Context) context.Context {
		q := Queue{}
		d := Duty{}
		d["default"] = "ok"
		q.Add("dynamic", d)
		t.Log("add d")
		t.Log(q)
		d2 := Duty{}
		d2["OK"] = "complete"
		q.Add("ok", d2)

		ctx = context.WithValue(ctx, "queue", q)
		return ctx
	}

	okFn := func(ctx context.Context) (context.Context, string, error) {
		ctx = context.WithValue(ctx, "key", "value")
		return ctx, "OK", nil
	}
	completeFn := func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, "complete", "finish")
		return ctx
	}

	wf := New()

	NewAgent(wf.Context, starterFn, "start")
	NewAgent(wf.Context, dynamicFn, "dynamic")
	NewAgent(wf.Context, completeFn, "complete")
	NewAgent(wf.Context, okFn, "ok")

	err = wf.Load(testConf)
	if err != nil {
		t.Log(testConf)
		t.Error(err)
	}
	err = wf.Handle("start")
	if err != nil {
		t.Error(err)
	}
	value, ok := wf.Context.Value("key").(string)
	t.Log(value, ok)
	if !ok {
		t.Error("value is not string")
	}
	if value != "value" {
		t.Error("value is not value", value)
	}

	value, ok = wf.Context.Value("complete").(string)
	t.Log(value, ok)
	if !ok {
		t.Error("complete is not string")
	}
	if value != "finish" {
		t.Error("complete is not finish", value)
	}
}
