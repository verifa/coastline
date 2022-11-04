package worker

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/verifa/coastline/requests"
	"github.com/verifa/coastline/server/oapi"
)

func TestWorker(t *testing.T) {
	opts := natsserver.DefaultTestOptions
	ns, err := server.NewServer(&opts)
	require.NoError(t, err)

	go ns.Start()
	if !ns.ReadyForConnections(time.Second) {
		t.Log("nats server could not start")
		t.FailNow()
	}

	engine, err := requests.Load(&requests.Config{
		Dir: "../requests/testdata",
	})
	require.NoError(t, err)

	{
		err := Start(engine, &Config{})
		require.NoError(t, err)
	}

	nc, err := nats.Connect(ns.ClientURL())
	require.NoError(t, err)

	tests := []struct {
		req       oapi.Request
		expectErr bool
	}{
		{
			req: oapi.Request{
				Kind: "t1",
				Spec: map[string]interface{}{
					"name": "project-x",
				},
			},
		},
		{
			req: oapi.Request{
				Kind: "t1",
				Spec: map[string]interface{}{
					"name":  "project-x",
					"error": "should not exist",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.req.Kind, func(t *testing.T) {
			done := make(chan struct{})

			sub, err := nc.Subscribe(subjectTaskResponse, func(msg *nats.Msg) {
				// This currently doesn't handle multiple tasks, i.e. if we receive multiple
				// messages because there are multiple tasks for a request, then we cannot
				// double-close a channel and will panic!
				close(done)
			})
			require.NoError(t, err)

			msg := TriggerMsg{
				Request: &tt.req,
			}
			reqBytes, err := json.Marshal(msg)
			require.NoError(t, err)
			{
				err := nc.Publish(subjectTriggerRun, reqBytes)
				require.NoError(t, err)
			}

			// Wait either until the channel "done" is closed, or a timeout is reached,
			// in which case fail the test case
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				assert.Fail(t, "timeout waiting for response")
			}

			// Unsubscribe from the queue so that we don't have multiple subscriptions
			{
				err := sub.Unsubscribe()
				require.NoError(t, err)
			}
		})
	}

}
