package requests

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/verifa/coastline/server/oapi"
)

type Template struct {
	Def   TemplateDef
	Value cue.Value
}

// TemplateDef defines the fields of a CUE-based Template for decoding
// and identifying which definitions in CUE are Request Templates
type TemplateDef struct {
	Kind    string `json:"kind"`
	Service struct {
		Selector struct {
			MatchLabels map[string]string `json:"matchLabels"`
		} `json:"selector"`
	} `json:"service"`
}

func (e *Engine) TemplatesForService(service *oapi.Service) []*Template {
	var templates []*Template
	for _, tmpl := range e.templates {
		for key, reqLabel := range tmpl.Def.Service.Selector.MatchLabels {
			serviceLabel, ok := service.Labels.Get(key)
			if ok && serviceLabel == reqLabel {
				templates = append(templates, tmpl)
			}
		}
	}
	return templates
}

// templateByKind returns the cue.Value of the request template for the
// given kind
func (e *Engine) templateByKind(kind string) (cue.Value, error) {
	for _, tmpl := range e.templates {
		if tmpl.Def.Kind == kind {
			return tmpl.Value, nil
		}
	}
	return cue.Value{}, fmt.Errorf("request template kind not found: %s", kind)
}

// getTemplates gets the request templates from the parsed cue files
func getTemplates(value cue.Value) ([]*Template, error) {
	var templates []*Template
	iter, err := value.Fields(cue.Definitions(true))
	if err != nil {
		return nil, fmt.Errorf("getting fields from cue value: %w", err)
	}
	kindPath := cue.ParsePath("kind")
	for iter.Next() {
		val := iter.Value()
		// Get the type value
		value := val.LookupPath(kindPath)
		if !value.Exists() {
			continue
		}
		var def TemplateDef
		if err := val.Decode(&def); err != nil {
			// Ignore errors, for now, and continue
			continue
		}
		if def.Kind != "" {
			templates = append(templates, &Template{
				Def:   def,
				Value: val,
			})
		}
	}
	return templates, nil
}
