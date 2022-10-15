package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjects(t *testing.T) {
	ctx := context.TODO()
	store, err := New(ctx)
	require.NoError(t, err)

	projects, err := store.Projects()
	require.NoError(t, err)

	fmt.Println(projects)
}
