package store

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/verifa/coastline/worker"
)

const (
	subjectTriggerRun   = "trigger.run"
	subjectTaskResponse = "task.response"
	queueServer         = "server"
)

func setupNATS() (*nats.Conn, error) {
	// Connect to the NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to nats: %w", err)
	}

	return nc, nil
}

func (s *Store) natsSubscribe() error {
	subs := []struct {
		subj    string
		handler nats.MsgHandler
	}{
		{
			subj:    subjectTaskResponse,
			handler: s.handleTaskResponse,
		},
	}
	for _, sub := range subs {
		_, err := s.nc.QueueSubscribe(sub.subj, queueServer, sub.handler)
		if err != nil {
			return fmt.Errorf("subscibing to queue %s: %w", queueServer, err)
		}
	}
	return nil
}

func (s *Store) handleTaskResponse(msg *nats.Msg) {
	var resp worker.ResponseMsg
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		log.Fatalln("unmarshalling task response: ", err.Error())
	}

	fmt.Println("Resp error: ", resp.Error)

	dbTask, err := s.client.Task.Create().
		SetTriggerID(resp.TriggerID).
		SetOutput(resp.Output).
		SetError(resp.Error).
		Save(s.ctx)
	if err != nil {
		log.Fatalln("saving task response: ", err.Error())
	}

	fmt.Printf("Saved task %s: %s\n", dbTask.ID.String(), string(dbTask.Output))
}

// IsNotFound returns a boolean indicating whether the error is a not found error.
// func IsNotFound(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	var e *NotFoundError
// 	return errors.As(err, &e)
// }
