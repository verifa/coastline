package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestsEngine(t *testing.T) {
	e, err := NewRequestsEngine(&RequestsEngineConfig{
		Module: "github.com/verifa/coastline/examples/basic",
	})
	require.NoError(t, err)
	fmt.Println(e)
}
