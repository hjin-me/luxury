package luxury

import (
	"testing"

	"golang.org/x/net/context"
)

func TestAgent(t *testing.T) {
	counter := 0

	ok := func(ctx context.Context) (context.Context, string, error) {
		counter = counter + 1
		return ctx, "OK", nil
	}
	wf := New()
	str := []string{"start", "step1", "step2", "step3", "step4"}
	NewAgent(wf.Context, ok, str...)
	for _, v := range str {
		agent, ok := GetAgent(wf.Context, v)
		if !ok {
			t.Error("agent not exsits")
		}
		_, status, err := agent(wf.Context)
		if err != nil {
			t.Error(err)
		}
		if status != "OK" {
			t.Error("status not ok", status)
		}
	}
	if counter != 5 {
		t.Error("counter not 5", counter)
	}
}

func TestAgentNoErrorNoCtx(t *testing.T) {

	counter := 0

	ok := func() string {
		counter = counter + 1
		return "OK"
	}
	wf := New()
	str := []string{"start", "step1", "step2", "step3", "step4"}
	NewAgent(wf.Context, ok, str...)
	for _, v := range str {
		agent, ok := GetAgent(wf.Context, v)
		if !ok {
			t.Error("agent not exsits")
		}
		_, status, err := agent(wf.Context)
		if err != nil {
			t.Error(err)
		}
		if status != "OK" {
			t.Error("status not ok", status)
		}
	}
	if counter != 5 {
		t.Error("counter not 5", counter)
	}
}
