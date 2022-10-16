package server

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode/gocodec"
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
	codec *gocodec.Codec
	value cue.Value
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

	return &RequestsEngine{
		codec: codec,
		value: requestDef,
	}, nil
}

func (e *RequestsEngine) Validate(input oapi.NewRequest) error {
	return e.codec.Validate(e.value, input)
}
