package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPolicies(t *testing.T) {
	pe := NewPolicyEngine()
	allow, err := pe.EvaluateLoginRequest(UserInfo{})
	require.NoError(t, err)
	t.Log("allow: ", allow)
	// TODO: fix test conditions
	// assert.True(t, allow)

}
