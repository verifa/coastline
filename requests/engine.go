package requests

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode/gocodec"
	"cuelang.org/go/encoding/openapi"
	"github.com/verifa/coastline/server/oapi"
)

type Config struct {
	Args []string
	Dir  string
}

func Load(config *Config) (*Engine, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	cuectx := cuecontext.New()
	buildInstances := load.Instances(config.Args, &load.Config{
		Dir: config.Dir,
	})
	if len(buildInstances) != 1 {
		return nil, fmt.Errorf("expecting only 1 build instance, got %d", len(buildInstances))
	}
	buildInstance := buildInstances[0]
	if buildInstance.Err != nil {
		return nil, fmt.Errorf("loading instance: %w", buildInstance.Err)
	}
	value := cuectx.BuildInstance(buildInstance)
	if value.Err() != nil {
		return nil, fmt.Errorf("building instance: %w", value.Err())
	}

	templates, err := getTemplates(value)
	if err != nil {
		return nil, fmt.Errorf("getting templates: %w", err)
	}
	tasks, err := getTasks(value)
	if err != nil {
		return nil, fmt.Errorf("getting tasks: %w", err)
	}
	codec := gocodec.New((*cue.Runtime)(cuectx), nil)

	return &Engine{
		cuectx:    cuectx,
		templates: templates,
		codec:     codec,
		tasks:     tasks,
		// value:     value,
	}, nil
}

// Engine is responsible for storing all the templates, validating incoming
// requests, serving OpenAPI specifications, and running worker cue-based tasks
type Engine struct {
	cuectx    *cue.Context
	codec     *gocodec.Codec
	templates []*Template
	tasks     []cue.Value
	// value     cue.Value
}

func (e *Engine) Validate(input oapi.NewRequest) error {
	cueInput := map[string]interface{}{
		"kind": input.Kind,
		"spec": input.Spec,
	}
	v, err := e.templateByKind(input.Kind)
	if err != nil {
		return err
	}
	return e.codec.Validate(v, cueInput)
}

// OpenAPISpec takes a Request Template kind and returns an OpenAPI JSON
// specification for it. If the Request Template kind does not exist, or there's an
// error generating the OpenAPI spec, an error is returned.
//
// There are little options to filter which Request Types should be included in
// the OpenAPI specification. What happens here is that an empty cue Instance is
// built and the specific Request Template is programmtically filled in to the
// empty instance
func (e *Engine) OpenAPISpec(kind string) ([]byte, error) {
	v, err := e.templateByKind(kind)
	if err != nil {
		return nil, fmt.Errorf("finding request type: %w", err)
	}

	// Build an empty file which we will fill with our request template
	val := e.cuectx.BuildFile(&ast.File{})
	if val.Err() != nil {
		return nil, fmt.Errorf("building empty file: %w", val.Err())
	}

	// Extract the spec field from the request template
	specValue := v.LookupPath(cue.ParsePath("spec"))
	if specValue.Err() != nil {
		return nil, fmt.Errorf("extracting spec field from request template: %w", specValue.Err())
	}

	// Make the path a definition by prefixing "#"
	fillPath := cue.ParsePath("#" + kind)
	reqTemplVal := val.FillPath(fillPath, specValue)
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
