package store

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/verifa/coastline/worker"
)

const (
	subjectTriggerRun       = "trigger.run"
	subjectWorkflowResponse = "workflow.response"
	queueServer             = "server"
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
			subj:    subjectWorkflowResponse,
			handler: s.handleWorkflowResponse,
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

func (s *Store) handleWorkflowResponse(msg *nats.Msg) {
	var resp worker.ResponseMsg
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		fmt.Println("Error: unmarshalling workflow response: ", err.Error())
	}

	fmt.Println("Resp error: ", resp.Error)
	var output map[string]interface{}
	if err := json.Unmarshal(resp.Output, &output); err != nil {
		fmt.Println("Error: unmarshalling response output: ", err.Error())
	}

	dbWorkflow, err := s.client.Workflow.Create().
		SetTriggerID(resp.TriggerID).
		SetOutput(output).
		SetError(resp.Error).
		Save(s.ctx)
	if err != nil {
		log.Fatalln("saving workflow response: ", err.Error())
	}

	fmt.Printf("Saved workflow %s: %s\n", dbWorkflow.ID.String(), dbWorkflow.Output)
}
