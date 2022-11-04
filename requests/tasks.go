package requests

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/verifa/coastline/server/oapi"
)

func (e *Engine) GetTasksForRequest(req *oapi.Request) ([]cue.Value, error) {
	var reqTasks = make([]cue.Value, 0)
	for _, task := range e.tasks {
		kind := task.LookupPath(cue.ParsePath("input.kind"))
		kindStr, err := kind.String()
		if err != nil {
			// This should not happen as list of tasks in engine should
			// already check this, but just to be sure
			return nil, fmt.Errorf("kind is not string: %w", err)
		}
		if kindStr == req.Kind {
			reqTasks = append(reqTasks, task)
		}
	}
	return reqTasks, nil
}

func getTasks(value cue.Value) ([]cue.Value, error) {
	taskPath := cue.ParsePath("task")
	taskVal := value.LookupPath(taskPath)
	if !taskVal.Exists() {
		return nil, nil
	}
	if taskVal.Kind() != cue.StructKind {
		return nil, fmt.Errorf("task must be struct")
	}

	taskIt, err := taskVal.Fields()
	if err != nil {
		return nil, fmt.Errorf("getting fields for task path: %w", err)
	}

	var tasks = make([]cue.Value, 0)
	for taskIt.Next() {
		kindPath := cue.ParsePath("input.kind")
		task := taskIt.Value()
		kindVal := task.LookupPath(kindPath)
		// If type doesn't exist or is not a string
		if !kindVal.Exists() || kindVal.Kind() != cue.StringKind {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func isTask(v cue.Value) bool {
	if v.Kind() != cue.StructKind {
		return false
	}
	if v.LookupPath(cue.ParsePath("input")).Exists() {
		return true
	}
	return false
}
