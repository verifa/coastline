package requests

import (
	"fmt"
	"io"

	"cuelang.org/go/cue"
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

	kindVal := v.LookupPath(cue.ParsePath("kind"))
	if !kindVal.Exists() {
		return nil, fmt.Errorf("request does not have a kind")
	}
	kindStr, err := kindVal.String()
	if err != nil {
		return nil, fmt.Errorf("request kind is not string: %w", err)
	}

	tmpl, err := e.templateByKind(kindStr)
	if err != nil {
		return nil, fmt.Errorf("getting template kind: %w", err)
	}

	// Requests might define some defaults, like the description, so we
	// need to Unify the values.
	// However, the NewRequest type contains fields like project_id which
	// are not accepted by default, so we create an "Accept" value from
	// the NewRequest type which allows these fields
	acceptVal := e.cuectx.EncodeType(oapi.NewRequest{})
	uniVal := tmpl.UnifyAccept(v, acceptVal)
	if uniVal.Err() != nil {
		return nil, fmt.Errorf("unifying new request with request template: %w", uniVal.Err())
	}
	// Decode the unified value into a NewRequest
	var req oapi.NewRequest
	if err := uniVal.Decode(&req); err != nil {
		return nil, fmt.Errorf("decoding cue value to new request: %w", err)
	}

	// Finally validate the input new request
	cueInput := map[string]interface{}{
		"kind": req.Kind,
		"spec": req.Spec,
	}
	if err := e.codec.Validate(tmpl, cueInput); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return &req, nil
}
