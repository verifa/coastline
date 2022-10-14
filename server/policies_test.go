package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicies(t *testing.T) {
	pe := NewPolicyEngine()
	allow, err := pe.EvaluateLoginRequest(UserInfo{})
	require.NoError(t, err)
	assert.True(t, allow)

}
