package requests

import (
	"fmt"
	"io"

	cuejson "cuelang.org/go/encoding/json"
	"github.com/verifa/coastline/server/oapi"
)

func (e *Engine) ValidateRequest(body io.Reader) (*oapi.NewRequest, error) {
	expr, err := cuejson.NewDecoder(nil, "", body).Extract()
	if err != nil {
		return nil, fmt.Errorf("extracting cue value from body: %w", err)
	}
	v := e.cuectx.BuildExpr(expr)
	if v.Err() != nil {
		return nil, fmt.Errorf("building cue value from body: %w", v.Err())
	}
	var req oapi.NewRequest
	if err := v.Decode(&req); err != nil {
		return nil, fmt.Errorf("decoding cue value to new request: %w", err)
	}

	tmpl, err := e.templateByKind(req.Kind)
	if err != nil {
		return nil, fmt.Errorf("getting template kind: %w", err)
	}
	cueInput := map[string]interface{}{
		"kind": req.Kind,
		"spec": req.Spec,
	}
	if err := e.codec.Validate(tmpl, cueInput); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return &req, nil
}
