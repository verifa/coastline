package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
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

	dummyUser := &oapi.User{
		Sub:    "123",
		Iss:    "123",
		Name:   "Testy Test",
		Groups: []string{"user"},
	}

	{
		// Simple case
		sessionID, err := store.NewSession(dummyUser)
		require.NoError(t, err)

		sessionUser, err := store.GetSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, dummyUser, sessionUser)
	}
	{
		// Same user, same client_id
		sessionID, err := store.NewSession(dummyUser)
		require.NoError(t, err)

		sessionUser, err := store.GetSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, dummyUser, sessionUser)
	}
	{
		anotherUser := *dummyUser
		anotherUser.Iss = "xyz"
		// Create under another client ID
		sessionID, err := store.NewSession(&anotherUser)
		require.NoError(t, err)

		sessionUser, err := store.GetSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, &anotherUser, sessionUser)
	}

	// Let's add and remove a group to the user
	dummyUser.Groups = []string{"admin"}
	{
		// Simple case
		sessionID, err := store.NewSession(dummyUser)
		require.NoError(t, err)

		sessionUser, err := store.GetSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, dummyUser, sessionUser)
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

	dummyUser := &oapi.User{
		Sub:    "123",
		Iss:    "123",
		Name:   "Testy Test",
		Groups: []string{"user"},
	}

	{
		sessionID, err := store.NewSession(dummyUser)
		require.NoError(t, err)

		_, err = store.GetSession(sessionID)
		require.NoError(t, err)

		time.Sleep(time.Second)

		_, err = store.GetSession(sessionID)
		// Expect session expired error
		assert.Error(t, err)
	}
}
