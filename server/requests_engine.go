package server

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode/gocodec"
	"cuelang.org/go/encoding/openapi"
	"github.com/verifa/coastline/server/oapi"
)

// RequestTemplate defines the fields of a CUE-based RequestTemplate for decoding
// and identifying which definitions in CUE are Request Templates
type RequestTemplate struct {
	Type    string `json:"type"`
	Service struct {
		Selector struct {
			MatchLabels map[string]string `json:"matchLabels"`
		} `json:"selector"`
	} `json:"service"`
}

type RequestsEngineConfig struct {
	Templates string
	Dir       string
}

// RequestsEngine is responsible for storing all the requests, validting incoming
// requets, and providing the frontend with request specs
type RequestsEngine struct {
	runtime  *cue.Runtime
	codec    *gocodec.Codec
	requests []*RequestTemplate
	instance *cue.Instance
}

func NewRequestsEngine(config *RequestsEngineConfig) (*RequestsEngine, error) {
	if config.Templates == "" {
		return nil, fmt.Errorf("templates required")
	}
	r := &cue.Runtime{}
	buildInstances := load.Instances([]string{config.Templates}, &load.Config{
		Dir: config.Dir,
	})
	if len(buildInstances) != 1 {
		return nil, fmt.Errorf("expecting only 1 build instance, got %d", len(buildInstances))
	}
	buildInstance := buildInstances[0]
	if buildInstance.Err != nil {
		return nil, fmt.Errorf("loading instance: %w", buildInstance.Err)
	}
	instance, err := r.Build(buildInstance)
	if err != nil {
		return nil, fmt.Errorf("building instance: %w", err)
	}
	if instance.Err != nil {
		return nil, fmt.Errorf("building instance: %w", buildInstance.Err)
	}

	requests, err := extractRequestTemplates(instance.Value())
	if err != nil {
		return nil, fmt.Errorf("extracting request definitions: %w", err)
	}
	codec := gocodec.New(r, nil)

	return &RequestsEngine{
		runtime:  r,
		requests: requests,
		codec:    codec,
		instance: instance,
	}, nil
}

func (e *RequestsEngine) Validate(input oapi.NewRequest) error {
	cueInput := map[string]interface{}{
		"type": input.Type,
		"spec": input.Spec,
	}
	v, err := e.requestTemplateByType(input.Type)
	if err != nil {
		return err
	}
	return e.codec.Validate(v, cueInput)
}

// OpenAPISpec takes the name of a Request Template and returns an OpenAPI
// specification for it. If the Request Template does not exist, or there's an
// error generating the OpenAPI spec, an error is returned.
//
// There are little options to filter which Request Types should be included in
// the OpenAPI specification. What happens here is that an empty cue Instance is
// built and the specific Request Template is programmtically filled in to the
// empty instance
func (e *RequestsEngine) OpenAPISpec(reqType string) ([]byte, error) {
	v, err := e.requestTemplateByType(reqType)
	if err != nil {
		return nil, fmt.Errorf("finding request type: %w", err)
	}

	// Build an empty instance which we will fill with our request template
	inst, err := e.runtime.CompileFile(&ast.File{})
	if err != nil {
		return nil, fmt.Errorf("compiling empty file: %w", err)
	}

	// Extract the spec field from the request template
	specValue := v.LookupPath(cue.ParsePath("spec"))
	if specValue.Err() != nil {
		return nil, fmt.Errorf("extracting spec field from request template: %w", specValue.Err())
	}

	// Make the path a definition by prefixing "#"
	fillPath := cue.ParsePath("#" + reqType)
	reqTemplVal := inst.Value().FillPath(fillPath, specValue)
	if reqTemplVal.Err() != nil {
		return nil, fmt.Errorf("error filling: %w", reqTemplVal.Err())
	}

	b, err := openapi.Gen(reqTemplVal, &openapi.Config{
		ExpandReferences: true,
	})
	if err != nil {
		return nil, fmt.Errorf("generating OpenAPI specification: %w", err)
	}

	return b, nil
}

func (e *RequestsEngine) RequestTemplatesForService(service *oapi.Service) []*RequestTemplate {
	var requests []*RequestTemplate
	for _, req := range e.requests {
		for key, reqLabel := range req.Service.Selector.MatchLabels {
			serviceLabel, ok := service.Labels.Get(key)
			if ok && serviceLabel == reqLabel {
				requests = append(requests, req)
			}
		}
	}
	return requests
}

// requestTemplateByType returns the cue.Value of the request template for the
// given type
func (e *RequestsEngine) requestTemplateByType(reqType string) (cue.Value, error) {
	iter, err := e.instance.Value().Fields(cue.All())
	if err != nil {
		return cue.Value{}, fmt.Errorf("getting fields from cue value: %w", err)
	}
	for iter.Next() {
		// Get the type value
		// reqPath := cue.MakePath()
		value := iter.Value().LookupPath(cue.ParsePath("type"))
		if !value.Exists() {
			continue
		}
		// Check the type value matches the given request type, and if so,
		// return the value
		if s, err := value.String(); err == nil {
			if s == reqType {
				return iter.Value(), nil
			}
		}
	}
	return cue.Value{}, fmt.Errorf("request type not found: %s", reqType)
}

// extractRequestTemplates gets the request templates from the parsed cue files
func extractRequestTemplates(value cue.Value) ([]*RequestTemplate, error) {
	var requests []*RequestTemplate
	iter, err := value.Fields(cue.All())
	if err != nil {
		return nil, fmt.Errorf("getting fields from cue value: %w", err)
	}
	for iter.Next() {
		var r RequestTemplate
		if err := iter.Value().Decode(&r); err != nil {
			// Ignore errors, for now, and continue
			continue
		}
		requests = append(requests, &r)
	}
	return requests, nil
}
