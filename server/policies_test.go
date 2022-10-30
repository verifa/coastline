package server

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestPolicies(t *testing.T) {
	pe := NewPolicyEngine()
	allow, err := pe.EvaluateLoginRequest(oapi.UserInfo{})
	require.NoError(t, err)
	t.Log("allow: ", allow)
	// TODO: fix test conditions
	// assert.True(t, allow)

}
