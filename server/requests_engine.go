package server

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode/gocodec"
	"cuelang.org/go/encoding/openapi"
	"github.com/verifa/coastline/server/oapi"
)

func DefaultRequestsEngineConfig() RequestsEngineConfig {
	return RequestsEngineConfig{
		Module:     "github.com/verifa/coastline/examples/basic",
		ModuleRoot: "./examples/basic",
	}
}

type RequestsEngineConfig struct {
	Module     string
	ModuleRoot string
}

// RequestsEngine is responsible for storing all the requests, validting incoming
// requets, and providing the frontend with request specs
type RequestsEngine struct {
	codec    *gocodec.Codec
	value    cue.Value
	instance *cue.Instance
	oapiSpec []byte
}

func NewRequestsEngine(config *RequestsEngineConfig) (*RequestsEngine, error) {
	r := &cue.Runtime{}
	buildInstances := load.Instances([]string{config.Module}, &load.Config{
		ModuleRoot: config.ModuleRoot,
		Module:     config.Module,
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

	// Get the Request definition from the provided cue model, which we will
	// use for validation
	cueRequestPath := cue.MakePath(cue.Def("Request"))
	requestDef := instance.Value().LookupPath(cueRequestPath)

	codec := gocodec.New(r, nil)

	// Generate OpenAPI Spec
	oapiSpec, err := openapi.Gen(instance, &openapi.Config{
		ExpandReferences: true,
		FieldFilter:      "Request",
	})

	if err != nil {
		return nil, fmt.Errorf("generating OpenAPI spec from instance: %w", err)
	}

	return &RequestsEngine{
		codec:    codec,
		instance: instance,
		value:    requestDef,
		oapiSpec: oapiSpec,
	}, nil
}

func (e *RequestsEngine) Validate(input oapi.NewRequest) error {
	cueInput := map[string]interface{}{
		"type": input.Type,
		"spec": input.Spec,
	}
	return e.codec.Validate(e.value, cueInput)
}

func (e *RequestsEngine) OpenAPISpec() []byte {
	return e.oapiSpec
}

// NOTE: this was a WIP as alternative to using OpenAPI with frontend...
// But OpenAPI would be MUCH easier/better
// func requestDefinitions(request cue.Value) (map[string]cue.Value, error) {
// 	var defs []cue.Value
// 	if request.IsConcrete() {
// 		if _, err := request.Struct(); err != nil {
// 			return nil, fmt.Errorf("#Request definition is concrete but not a struct")
// 		}
// 		// If request is a concrete struct, then only a single request type exists
// 		// and we will use it. Likely the user is testing because in reality we
// 		// want more than one value
// 		defs = append(defs, request)
// 	} else {
// 		op, values := request.Expr()
// 		if op != cue.OrOp {
// 			return nil, fmt.Errorf("#Request definition must be a struct or multiple structs using the \"|\" operator")
// 		}
// 		defs = values
// 	}

// 	typeMap := make(map[string]cue.Value)
// 	// Convert the values into a map based on the "Type" field
// 	for _, def := range defs {
// 		// There should be a better way of doing this with LookupPath,
// 		// but it didn't work easily so let's iterate over the fields and find
// 		// the "type"
// 		it, err := def.Fields()
// 		if err != nil {
// 			return nil, fmt.Errorf("cannot get fields for value: %w", err)
// 		}
// 		var (
// 			hasType  bool
// 			typeName string
// 		)
// 		for it.Next() && !hasType {
// 			if it.Selector().String() == "type" {
// 				hasType = true
// 				name, err := it.Value().String()
// 				if err != nil {
// 					return nil, fmt.Errorf("\"type\" field must be a string: %w", err)
// 				}
// 				typeName = name
// 			}
// 		}
// 		if !hasType {
// 			return nil, fmt.Errorf("#Request definition must include a type")
// 		}

// 		typeMap[typeName] = def
// 	}
// 	return typeMap, nil
// }
