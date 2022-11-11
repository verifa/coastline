package requests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestRunner(t *testing.T) {
	engine, err := Load(&Config{
		Dir: "./testdata",
	})
	require.NoError(t, err)

	req := oapi.Request{
		Kind: "t1",
		Spec: map[string]interface{}{
			"name": "test",
		},
	}

	paths, err := engine.GetWorkflowsForRequest(&req)
	require.NoError(t, err)

	for _, path := range paths {
		v, err := engine.RunWorkflow(path, &req)
		require.NoError(t, err)
		b, err := v.MarshalJSON()
		require.NoError(t, err)
		fmt.Println(string(b))
	}
}
