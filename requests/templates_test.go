package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestTemplatesByKind(t *testing.T) {
	config := Config{
		Dir: "./testdata",
	}
	e, err := Load(&config)
	require.NoError(t, err)

	tests := []struct {
		kind      string
		expectErr bool
	}{
		{
			kind: "t1",
		},
		{
			kind: "t2",
		},
		{
			kind:      "unknown",
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			tmpl, err := e.templateByKind(tt.kind)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tmpl)
			}
		})
	}
}

func TestTemplatesForService(t *testing.T) {
	config := Config{
		Dir: "./testdata",
	}
	e, err := Load(&config)
	require.NoError(t, err)

	tests := []struct {
		service   oapi.Service
		expectErr bool
	}{
		{
			service: oapi.Service{
				Name: "s1",
				Labels: &oapi.Service_Labels{
					AdditionalProperties: map[string]string{
						"tool": "t1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.service.Name, func(t *testing.T) {
			tmpls := e.TemplatesForService(&tt.service)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, tmpls)
			}
		})
	}
}
