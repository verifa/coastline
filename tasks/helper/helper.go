package helper

import (
	"fmt"
	"io"

	"cuelang.org/go/cue"
	"cuelang.org/go/tools/flow"
	"github.com/hashicorp/go-multierror"
)

type TaskHelper struct {
	Task *flow.Task
	Errs error
}

func (h *TaskHelper) MustDecode(path string, v interface{}) {
	p := h.Task.Value().LookupPath(cue.ParsePath(path))
	err := p.Decode(v)
	if err != nil {
		h.AppendErr(fmt.Errorf("decoding value at path %s: %w", path, err))
	}
}

func (h *TaskHelper) Decode(path string, v interface{}) bool {
	p := h.Task.Value().LookupPath(cue.ParsePath(path))
	if !p.Exists() {
		return false
	}
	h.MustDecode(path, v)
	return true
}

func (h *TaskHelper) MustString(path string) string {
	p := h.Task.Value().LookupPath(cue.ParsePath(path))
	value, err := p.String()
	if err != nil {
		h.AppendErr(fmt.Errorf("invalid string at path %s: %w", path, err))
		return ""
	}
	return value
}

func (h *TaskHelper) String(path string) (string, bool) {
	p := h.Task.Value().LookupPath(cue.ParsePath(path))
	if !p.Exists() {
		return "", false
	}
	return h.MustString(path), true
}

func (h *TaskHelper) MustReader(path string) io.Reader {
	p := h.Task.Value().LookupPath(cue.ParsePath(path))
	value, err := p.Reader()
	if err != nil {
		h.AppendErr(fmt.Errorf("invalid reader at path %s: %w", path, err))
		return nil
	}
	return value
}

func (h *TaskHelper) Reader(path string) (io.Reader, bool) {
	p := h.Task.Value().LookupPath(cue.ParsePath(path))
	if !p.Exists() {
		return nil, false
	}
	return h.MustReader(path), true
}

func (h *TaskHelper) AppendErr(err error) {
	h.Errs = multierror.Append(h.Errs, err)
}
