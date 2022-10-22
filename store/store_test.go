package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/server/oapi"
)

func TestStore(t *testing.T) {
	ctx := context.TODO()
	store, err := New(ctx)
	require.NoError(t, err)

	newProject := oapi.NewProject{
		Name: "MyProject",
	}
	project, err := store.CreateProject(&newProject)
	require.NoError(t, err)

	newService := oapi.NewService{
		Name: "MyService",
	}
	service, err := store.CreateService(&newService)
	require.NoError(t, err)

	newRequest := oapi.NewRequest{
		Type:        "test",
		RequestedBy: "me",
		ProjectId:   project.Id,
		ServiceId:   service.Id,
		Spec:        map[string]interface{}{"request_key": "request_value"},
	}
	request, err := store.CreateRequest(&newRequest)
	require.NoError(t, err)
	fmt.Println("request: ", request)
	spew.Dump(request)

	projectResp, err := store.QueryProjects()
	require.NoError(t, err)
	fmt.Println(projectResp)

	{
		review, err := store.CreateReview(request.Id, &oapi.NewReview{
			Status: oapi.NewReviewStatusApprove,
			Type:   oapi.NewReviewTypeUser,
		})
		require.NoError(t, err)
		spew.Dump(review)
	}
	{
		review, err := store.CreateReview(request.Id, &oapi.NewReview{
			Status: oapi.NewReviewStatusApprove,
			Type:   oapi.NewReviewTypeUser,
		})
		require.NoError(t, err)
		spew.Dump(review)
	}

	requestResp, err := store.QueryRequests()
	require.NoError(t, err)
	spew.Dump(requestResp)
}
