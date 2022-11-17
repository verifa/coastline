package requests

import (
	"context"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/tools/flow"
)

func (e *Engine) getTests(filter string) ([]cue.Path, error) {
	var tests = make([]cue.Path, 0)
	testPath := cue.ParsePath("test")
	testVal := e.value.LookupPath(testPath)
	// If there's no tests, it's fine, return the empty list
	if !testVal.Exists() {
		return tests, nil
	}
	testIt, err := testVal.Fields()
	if err != nil {
		return nil, fmt.Errorf("getting field iterator for path %s: %w", testPath, err)
	}
	for testIt.Next() {
		// TODO: validate that the test is a valid test, somehow...
		v := testIt.Value()
		if strings.HasPrefix(v.Path().String(), filter) {
			tests = append(tests, testIt.Value().Path())
		}
	}
	return tests, nil
}

type TestConfig struct {
	Filter string
}

func (e *Engine) RunTests(config *TestConfig) (*TestResult, error) {
	if config == nil {
		config = &TestConfig{}
	}
	tests, err := e.getTests(config.Filter)
	if err != nil {
		return nil, fmt.Errorf("getting tests: %w", err)
	}

	var t TestResult

	for _, test := range tests {
		t.Run(test, func(t *TestCase) {

			testVal := e.value.LookupPath(test)
			if !testVal.Exists() {
				t.Failf("test path does not exist: %s", test)
				return
			}

			testWf := testVal.LookupPath(cue.ParsePath("run"))
			if !testWf.Exists() {
				t.Failf("test does not specify a workflow to run: %s", test)
				return
			}
			_, wfPath := testWf.ReferencePath()
			if wfPath.Err() != nil {
				t.Failf("cannot get workflow reference at %s: %s", testWf.Path().String(), err.Error())
				return
			}

			r := testRunner{
				t:    t,
				path: test,
				wf: workflowRunner{
					path: wfPath,
				},
			}
			controller := flow.New(&flow.Config{
				Root:           r.path,
				IgnoreConcrete: true,
			}, e.value, r.workflowFunc)
			err := controller.Run(context.Background())
			if err != nil {
				t.Failf("cue controller error: %s", err.Error())
			}
			t.Output = controller.Value().LookupPath(test)
		})
	}

	return &t, nil
}

type testRunner struct {
	t    *TestCase
	path cue.Path
	wf   workflowRunner
}

func (r *testRunner) workflowFunc(v cue.Value) (flow.Runner, error) {
	fr, err := r.wf.workflowFunc(v)
	if fr != nil || err != nil {
		return fr, err
	}

	if isTestAssert(v) {
		at := assertTask{
			t: r.t,
		}

		return at, nil
	}
	return nil, nil
}

func isTestAssert(v cue.Value) bool {
	path := v.Path()
	if len(path.Selectors()) != 3 {
		return false
	}
	sel := path.Selectors()
	endSel := sel[len(sel)-1]

	return endSel.String() == "assert"
}

type assertTask struct {
	t *TestCase
}

func (at assertTask) Run(t *flow.Task, pErr error) error {
	assertIt, err := t.Value().Fields()
	if err != nil {
		return fmt.Errorf("getting assert fields: %w", err)
	}
	for assertIt.Next() {
		v := assertIt.Value()
		if v.Kind() != cue.BoolKind {
			return fmt.Errorf("assert must be bool: %s", v.Path())
		}
		vBool, err := v.Bool()
		if err != nil {
			return fmt.Errorf("getting bool value at %s: %w", v.Path(), err)
		}
		if !vBool {
			at.t.Failf("assert: expected true but was false at: %s", v.Path())
			continue
		}
	}
	return nil
}
