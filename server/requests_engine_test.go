package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestRequestsEngine(t *testing.T) {
	e, err := NewRequestsEngine(&RequestsEngineConfig{
		Module:     "github.com/verifa/coastline/examples/basic",
		ModuleRoot: "../examples/basic",
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

	// []oapi.NewRequest{
	// }
}
