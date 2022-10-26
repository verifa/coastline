package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestTemplatesForService(t *testing.T) {
	e, err := NewRequestsEngine(&RequestsEngineConfig{
		Dir:       "..",
		Templates: "./examples/basic",
	})
	require.NoError(t, err)

	reqs := e.RequestTemplatesForService(&oapi.Service{
		Labels: &oapi.Service_Labels{
			AdditionalProperties: map[string]string{
				"tool": "artifactory",
			},
		},
	})
	require.NoError(t, err)
	assert.Len(t, reqs, 1)

	b, err := e.OpenAPISpec("ArtifactoryRepoRequest")
	require.NoError(t, err)
	t.Log("spec: ", string(b))
}

func TestRequestsEngine(t *testing.T) {
	e, err := NewRequestsEngine(&RequestsEngineConfig{
		Dir:       "..",
		Templates: "./examples/basic",
	})
	require.NoError(t, err)

	tests := []struct {
		req       oapi.NewRequest
		expectErr bool
	}{
		{
			req: oapi.NewRequest{
				Type:        "ArtifactoryRepoRequest",
				RequestedBy: "someone",
				Spec: map[string]interface{}{
					"repo": "hello",
				},
			},
			expectErr: false,
		},
		{
			req: oapi.NewRequest{
				Type:        "JenkinsServerRequest",
				RequestedBy: "someone",
				Spec: map[string]interface{}{
					"name": "hello",
				},
			},
			expectErr: false,
		},
		{
			req: oapi.NewRequest{
				Type:        "InvalidJenkinsServerRequest",
				RequestedBy: "someone",
				Spec: map[string]interface{}{
					"repo": "hello",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.req.Type, func(t *testing.T) {
			err := e.Validate(tt.req)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
