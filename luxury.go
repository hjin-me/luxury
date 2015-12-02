package luxury

import (
	"errors"
	"sync/atomic"

	"github.com/hjin-me/luxury/logex"

	"github.com/hjin-me/luxury/config"

	"golang.org/x/net/context"
)

type Duty map[string]string
type Workflow struct {
	Context context.Context
	Queue   Queue
	Now     string
	Deep    uint32
}
type Queue map[string]Duty

func (q Queue) Add(name string, d Duty) {
	q[name] = d
}

func New() Workflow {
	wf := Workflow{}
	t := make(WorkflowAgent)
	wfa := &t
	wf.Context = context.WithValue(context.Background(), "wf", wfa)
	wf.Queue = make(map[string]Duty)
	return wf
}

func (wf *Workflow) Load(filename string) error {
	return config.Load(filename, &wf.Queue)
}

func (wf *Workflow) Add(name string, duty Duty) {
	wf.Queue[name] = duty
}

func (wf *Workflow) Handle(status string) error {
	wf.Deep = 0
	wf.Now = status
	return wf.next()
}
func (wf *Workflow) next() error {
	for {
		atomic.AddUint32(&wf.Deep, 1)
		agent, ok := GetAgent(wf.Context, wf.Now)
		if !ok {
			return errors.New("agent named [" + wf.Now + "] not found")
		}
		var (
			next string
			err  error
		)
		wf.Context, next, err = agent(wf.Context)
		logex.Trace("[switch]", wf.Now, next)
		if err != nil {
			return err
		}
		// add new step
		queue, ok := wf.Context.Value("queue").(Queue)
		if ok {
			for k, v := range queue {
				wf.Queue[k] = v
			}
			queue = nil
		}

		// agent process
		if next == "end" {
			return nil
		}

		// find next step
		d, ok := wf.Queue[wf.Now]
		if !ok {
			return nil
		}
		wf.Now, _ = d[next]
	}

	return nil

}
