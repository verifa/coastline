package requests

import (
	"context"
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cuejson "cuelang.org/go/encoding/json"
	"cuelang.org/go/tools/flow"
	"github.com/verifa/coastline/server/oapi"
	"github.com/verifa/coastline/tasks"
)

// RunWorkflow takes a cue path to a workflow and a request and runs the workflow at
// the given path, replacing the input with the spec from the given request
func (e *Engine) RunWorkflow(wfPath cue.Path, req *oapi.Request) (cue.Value, error) {
	r := workflowRunner{
		path: wfPath,
		req:  req,
	}
	wf := e.value.LookupPath(wfPath)
	if !wf.Exists() {
		return cue.Value{}, fmt.Errorf("workflow does not exist")
	} else if !isWorkflow(wf) {
		// TODO: we could be a bit more descriptive here
		return cue.Value{}, fmt.Errorf("workflow exists but is not valid")
	}

	fmt.Println("Running workflow: ", wfPath)
	controller := flow.New(&flow.Config{
		Root:           wfPath,
		IgnoreConcrete: true,
	}, e.value, r.workflowFunc)
	err := controller.Run(context.Background())
	if err != nil {
		return cue.Value{}, fmt.Errorf("running workflow: %w", err)
	}
	// Get final result and output
	result := controller.Value()
	outputSel := cue.ParsePath("output").Selectors()[0]
	outputPath := cue.MakePath(append(wfPath.Selectors(), outputSel)...)
	output := result.LookupPath(outputPath)
	// If output doesn't exist return an empty value
	if !output.Exists() {
		return cue.Value{}, nil
	}

	return output, nil
}

type workflowRunner struct {
	path cue.Path
	req  *oapi.Request
}

func (r *workflowRunner) workflowFunc(v cue.Value) (flow.Runner, error) {
	if isWorkflowInput(v) {
		it := inputTask{
			req: r.req,
		}
		return &it, nil
	}
	// Handle Coastline tasks
	if taskID, ok := tasks.TaskID(v); ok {
		ct := tasks.Task{
			ID: taskID,
		}
		return ct, nil
	}
	return nil, nil
}

var _ flow.Runner = (*inputTask)(nil)

type inputTask struct {
	req *oapi.Request
}

// Tasks must implement a Run func, this is where we execute our task
func (it inputTask) Run(t *flow.Task, pErr error) error {
	// Converting from JSON to Go is troublesome because of numbers.
	// Are they floats, or int8, 16, uint16, etc.
	// Cue's JSON package can be used to build a CUE value, so let's use that
	// instead.
	if it.req == nil {
		return nil
	}
	// First re-create JSON from the request spec
	specBytes, err := json.Marshal(it.req.Spec)
	if err != nil {
		return fmt.Errorf("marshalling request spec: %w", err)
	}

	// Validate the JSON spec against the spec value from the input request
	specVal := t.Value().LookupPath(cue.ParsePath("spec"))
	if !specVal.Exists() {
		return fmt.Errorf("spec does not exist in input")
	}
	if err := cuejson.Validate(specBytes, specVal); err != nil {
		return fmt.Errorf("spec is not valid: %w", err)
	}
	expr, err := cuejson.Extract("", specBytes)
	if err != nil {
		return fmt.Errorf("extracting cue value from json: %w", err)
	}

	v := cuecontext.New().BuildExpr(expr)
	return t.Fill(map[string]interface{}{
		"spec": v,
	})
}

func isWorkflowInput(v cue.Value) bool {
	// TODO: should this be stricter at checking?
	// pathSelectors := v.Path().Selectors()
	// curSelector := pathSelectors[len(pathSelectors)-1]
	// if curSelector.String() != "input" {
	// 	return false
	// }
	if v.Kind() != cue.StructKind {
		return false
	}
	if v.LookupPath(cue.ParsePath("kind")).Exists() {
		return true
	}
	return false
}
