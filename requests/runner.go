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
)

func NewRunner(req *oapi.Request) *runner {
	return &runner{
		req: req,
	}
}

type runner struct {
	req *oapi.Request
}

func (r *runner) RunTask(v cue.Value) (cue.Value, error) {
	if !isTask(v) {
		return cue.Value{}, fmt.Errorf("value is not task")
	}
	fmt.Println("RUNNING TASK: ", v.Path().String())
	controller := flow.New(&flow.Config{}, v, r.taskFunc)
	err := controller.Run(context.Background())
	if err != nil {
		return cue.Value{}, fmt.Errorf("running task: %w", err)
	}
	// Get final result and output
	result := controller.Value()
	outputPath := cue.ParsePath("output")
	output := result.LookupPath(outputPath)
	// If output doesn't exist return an empty value
	if !output.Exists() {
		return cue.Value{}, nil
	}

	return output, nil
}

func (r *runner) taskFunc(v cue.Value) (flow.Runner, error) {
	if isTaskInput(v) {
		return flow.RunnerFunc(func(t *flow.Task) error {
			// Converting from JSON to Go is troublesome because of numbers.
			// Are they floats, or int8, 16, uint16, etc.
			// Cue's JSON package can be used to build a CUE value, so let's use that
			// instead.
			// First re-create JSON from the request spec
			specBytes, err := json.Marshal(r.req.Spec)
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
		}), nil
	}
	// Skip all values that are not tasks
	if isBuiltinTask(v) {
		return flow.RunnerFunc(func(t *flow.Task) error {
			// How can I run this builtin task from here ??
			// E.g. https://pkg.go.dev/cuelang.org/go/pkg/tool
			return nil
		}), nil
	}
	return nil, nil
}

func isTaskInput(v cue.Value) bool {
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

func isBuiltinTask(v cue.Value) bool {
	if v.Kind() != cue.StructKind {
		return false
	}
	if v.LookupPath(cue.ParsePath("$id")).Exists() {
		return true
	}
	return false
}
