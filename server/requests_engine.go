package server

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode/gocodec"
)

type RequestsEngineConfig struct {
	Module     string
	ModuleRoot string
}

// RequestsEngine is responsible for storing all the requests, validting incoming
// requets, and providing the frontend with request specs
type RequestsEngine struct {
	codec    *gocodec.Codec
	instance *cue.Instance
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

	codec := gocodec.New(r, nil)

	return &RequestsEngine{
		codec:    codec,
		instance: instance,
	}, nil
}

func (e *RequestsEngine) Validate(input interface{}) error {
	return e.codec.Validate(e.instance.Value(), input)
}
