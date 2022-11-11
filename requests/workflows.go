package requests

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/verifa/coastline/server/oapi"
)

func (e *Engine) GetWorkflowsForRequest(req *oapi.Request) ([]cue.Path, error) {
	var reqWorkflows = make([]cue.Path, 0)
	for _, workflow := range e.workflows {
		kind := workflow.LookupPath(cue.ParsePath("input.kind"))
		kindStr, err := kind.String()
		if err != nil {
			// This should not happen as list of workflows in engine should
			// already check this, but just to be sure
			return nil, fmt.Errorf("kind is not string: %w", err)
		}
		if kindStr == req.Kind {
			reqWorkflows = append(reqWorkflows, workflow.Path())
		}
	}
	return reqWorkflows, nil
}

func getWorkflows(value cue.Value) ([]cue.Value, error) {
	workflowPath := cue.ParsePath("workflow")
	workflowVal := value.LookupPath(workflowPath)
	if !workflowVal.Exists() {
		return nil, nil
	}
	if workflowVal.Kind() != cue.StructKind {
		return nil, fmt.Errorf("workflow must be struct")
	}

	workflowIt, err := workflowVal.Fields()
	if err != nil {
		return nil, fmt.Errorf("getting fields for workflow path: %w", err)
	}

	var workflows = make([]cue.Value, 0)
	for workflowIt.Next() {
		kindPath := cue.ParsePath("input.kind")
		workflow := workflowIt.Value()
		kindVal := workflow.LookupPath(kindPath)
		// If type doesn't exist or is not a string
		if !kindVal.Exists() || kindVal.Kind() != cue.StringKind {
			continue
		}
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

func isWorkflow(v cue.Value) bool {
	if v.Kind() != cue.StructKind {
		return false
	}
	if v.LookupPath(cue.ParsePath("input")).Exists() {
		return true
	}
	return false
}
