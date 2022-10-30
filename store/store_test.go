package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestStore(t *testing.T) {
	ctx := context.TODO()
	store, err := New(ctx, &Config{})
	require.NoError(t, err)

	user := oapi.User{
		Sub:  "123",
		Iss:  "123",
		Name: "Bob",
	}

	{
		_, err := store.createUpdateUser(&user)
		require.NoError(t, err)
	}

	newProject := oapi.NewProject{
		Name: "MyProject",
	}
	project, err := store.CreateProject(&newProject)
	require.NoError(t, err)

	projectResp, err := store.QueryProjects()
	require.NoError(t, err)
	assert.Len(t, projectResp.Projects, 1)

	newService := oapi.NewService{
		Name: "MyService",
	}
	service, err := store.CreateService(&newService)
	require.NoError(t, err)

	serviceResp, err := store.QueryServices()
	require.NoError(t, err)
	assert.Len(t, serviceResp.Services, 1)

	newRequest := oapi.NewRequest{
		Type:      "test",
		ProjectId: project.Id,
		ServiceId: service.Id,
		Spec:      map[string]interface{}{"request_key": "request_value"},
	}
	request, err := store.CreateRequest(&user, &newRequest)
	require.NoError(t, err)

	{
		_, err := store.CreateReview(request.Id, &user, &oapi.NewReview{
			Status: oapi.NewReviewStatusApprove,
			Type:   oapi.NewReviewTypeUser,
		})
		require.NoError(t, err)
	}
	{
		_, err := store.CreateReview(request.Id, &user, &oapi.NewReview{
			Status: oapi.NewReviewStatusApprove,
			Type:   oapi.NewReviewTypeUser,
		})
		require.NoError(t, err)
	}

	requestResp, err := store.QueryRequests()
	require.NoError(t, err)
	assert.Len(t, requestResp.Requests, 1)
	assert.Len(t, requestResp.Requests[0].Reviews, 2)
}
