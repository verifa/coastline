package requests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunnerTest(t *testing.T) {
	e, err := Load(&Config{
		Dir:       "../examples/basic",
		IsTesting: true,
	})
	require.NoError(t, err)

	result, err := e.RunTests(nil)
	require.NoError(t, err)

	result.Print(false)
}
