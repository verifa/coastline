package requests

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestValidateRequest(t *testing.T) {
	config := Config{
		Dir: "./testdata",
	}
	e, err := Load(&config)
	require.NoError(t, err)

	tests := []struct {
		input     oapi.NewRequest
		expectErr bool
	}{
		{
			input: oapi.NewRequest{
				Kind: "t1",
				Spec: map[string]interface{}{
					"name": "hello",
				},
			},
		},
		{
			input: oapi.NewRequest{
				Kind: "t1",
				Spec: map[string]interface{}{
					"name":  "hello",
					"error": "should not exist",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input.Kind, func(t *testing.T) {
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tt.input)
			require.NoError(t, err)
			req, err := e.ValidateRequest(&buf)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.input.Kind, req.Kind)
				assert.Equal(t, tt.input.Spec, req.Spec)
			}
		})
	}

}
