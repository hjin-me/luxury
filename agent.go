package luxury

import "golang.org/x/net/context"

var (
	agentPool = make(map[string]formatAgent)
)

type WorkflowAgent map[string]formatAgent

func (wfa *WorkflowAgent) Set(fn formatAgent, alias ...string) {
	wf := *wfa
	for _, name := range alias {
		wf[name] = fn
	}
	wfa = &wf
}

func NewAgent(ctx context.Context, fn Agent, alias ...string) {
	ffn := wrap(fn)
	wfa, ok := ctx.Value("wf").(*WorkflowAgent)
	if !ok {
		panic("WorkflowAgent not found")
	}
	wfa.Set(ffn, alias...)
}

func GetAgent(ctx context.Context, name string) (formatAgent, bool) {
	wfa, ok := ctx.Value("wf").(*WorkflowAgent)
	if !ok {
		return nil, false
	}
	if ag, ok := (*wfa)[name]; ok {
		return ag, true
	}
	return nil, false
}
