package tasks

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/tools/flow"
	"github.com/verifa/coastline/tasks/http"
)

var _ flow.Runner = (*Task)(nil)

type Task struct {
	ID string
}

func (t Task) Run(task *flow.Task, pErr error) error {
	var runner flow.Runner
	switch t.ID {
	case "http.Do":
		runner = http.DoTask{}
	default:
		return fmt.Errorf("unknown task: %s", t.ID)

	}
	return runner.Run(task, pErr)
}

func TaskID(v cue.Value) (string, bool) {
	if v.Kind() != cue.StructKind {
		return "", false
	}
	taskVal := v.LookupPath(cue.ParsePath("$task"))
	if taskVal.Exists() && taskVal.Kind() == cue.StringKind {
		taskStr, _ := taskVal.String()
		return taskStr, true
	}
	return "", false
}
