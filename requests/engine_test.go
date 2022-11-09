package requests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestRequests(t *testing.T) {
	config := Config{
		Dir: "./testdata",
		// Args: []string{"./t1.cue"},
	}
	e, err := Load(&config)
	require.NoError(t, err)

	tests := []struct {
		req       oapi.NewRequest
		expectErr bool
	}{
		{
			req: oapi.NewRequest{
				Kind: "t1",
				Spec: map[string]interface{}{
					"name": "hello",
				},
			},
			expectErr: false,
		},
		{
			req: oapi.NewRequest{
				Kind: "t1",
				Spec: map[string]interface{}{
					"name": "hello",
				},
			},
			expectErr: false,
		},
		{
			req: oapi.NewRequest{
				Kind: "t1",
				Spec: map[string]interface{}{
					"name": "hello has space",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.req.Kind, func(t *testing.T) {
			err := e.Validate(tt.req)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOpenAPI(t *testing.T) {
	config := Config{
		Dir: "./testdata",
	}
	e, err := Load(&config)
	require.NoError(t, err)

	b, err := e.OpenAPISpec("t1")
	require.NoError(t, err)
	fmt.Println(string(b))

}
