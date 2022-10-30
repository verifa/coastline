package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/session"
)

func TestGroup(t *testing.T) {
	ctx := context.TODO()
	store, err := New(ctx, &Config{})
	require.NoError(t, err)

	testGroups := []string{"user", "admin", "team"}
	{
		entGroups, err := store.createReadGroups(testGroups)
		require.NoError(t, err)
		assert.Len(t, entGroups, len(testGroups))
	}
	testGroups = append(testGroups, "another")
	{
		entGroups, err := store.createReadGroups(testGroups)
		require.NoError(t, err)
		assert.Len(t, entGroups, len(testGroups))
	}

}

func TestSession(t *testing.T) {
	ctx := context.TODO()
	store, err := New(ctx, &Config{})
	require.NoError(t, err)

	claims := &session.UserClaims{
		Sub:    "123",
		Email:  "test@localhost",
		Name:   "Testy Test",
		Groups: []string{"user"},
	}

	{
		// Simple case
		sessionID, err := store.NewSession("abc", claims)
		require.NoError(t, err)

		retClaims, err := store.ValidateSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, claims, retClaims)
	}
	{
		// Same user, same client_id
		sessionID, err := store.NewSession("abc", claims)
		require.NoError(t, err)

		retClaims, err := store.ValidateSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, claims, retClaims)
	}
	{
		// Create under another client ID
		sessionID, err := store.NewSession("xyz", claims)
		require.NoError(t, err)

		retClaims, err := store.ValidateSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, claims, retClaims)
	}

	// Let's add and remove a group to the claims
	claims.Groups = []string{"admin"}
	{
		// Simple case
		sessionID, err := store.NewSession("abc", claims)
		require.NoError(t, err)

		retClaims, err := store.ValidateSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, claims, retClaims)
	}

	// Check we have 2 sessions and 2 users
	assert.Len(t, store.client.Session.Query().IDsX(ctx), 4)
	assert.Len(t, store.client.User.Query().IDsX(ctx), 2)
}

func TestSessionExpired(t *testing.T) {
	ctx := context.TODO()
	store, err := New(ctx, &Config{
		SessionDuration: time.Second,
	})
	require.NoError(t, err)

	claims := &session.UserClaims{
		Sub:    "123",
		Email:  "test@localhost",
		Name:   "Testy Test",
		Groups: []string{"user"},
	}

	{
		sessionID, err := store.NewSession("abc", claims)
		require.NoError(t, err)

		_, err = store.ValidateSession(sessionID)
		require.NoError(t, err)

		time.Sleep(time.Second)

		_, err = store.ValidateSession(sessionID)
		// Expect session expired error
		assert.Error(t, err)
	}
}
