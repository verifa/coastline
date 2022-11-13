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

type ServiceDef struct {
	Name string `json:"name"`
}

// TemplateDef defines the fields of a CUE-based Template for decoding
// and identifying which definitions in CUE are Request Templates
type TemplateDef struct {
	Kind     string       `json:"kind"`
	Services []ServiceDef `json:"services"`
	// Service struct {
	// 	Selector struct {
	// 		MatchLabels map[string]string `json:"matchLabels"`
	// 	} `json:"selector"`
	// } `json:"service"`
}

func (e *Engine) TemplatesForService(service *oapi.Service) []*Template {
	var templates []*Template
	for _, tmpl := range e.templates {
		for _, svc := range tmpl.Def.Services {
			if svc.Name == service.Name {
				templates = append(templates, tmpl)
			}
		}
	}
	// for _, tmpl := range e.templates {
	// 	for key, reqLabel := range tmpl.Def.Service.Selector.MatchLabels {
	// 		serviceLabel, ok := service.Labels.Get(key)
	// 		if ok && serviceLabel == reqLabel {
	// 			templates = append(templates, tmpl)
	// 		}
	// 	}
	// }
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
	templatesPath := cue.ParsePath("request")
	templatesVal := value.LookupPath(templatesPath)
	if !templatesVal.Exists() {
		return nil, nil
	}
	var templates []*Template
	iter, err := templatesVal.Fields(cue.Definitions(true))
	if err != nil {
		return nil, fmt.Errorf("getting fields from cue value: %w", err)
	}
	kindPath := cue.ParsePath("kind")
	for iter.Next() {
		tmplVal := iter.Value()
		// Get the type value
		kindVal := tmplVal.LookupPath(kindPath)
		if !kindVal.Exists() {
			continue
		}
		var def TemplateDef
		if err := tmplVal.Decode(&def); err != nil {
			// Ignore errors, for now, and continue
			continue
		}
		if def.Kind != "" {
			templates = append(templates, &Template{
				Def:   def,
				Value: tmplVal,
			})
		}
	}
	return templates, nil
}
