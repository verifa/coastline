package worker

import (
	"encoding/json"
	"fmt"
	"log"

	"cuelang.org/go/cue"
	"github.com/nats-io/nats.go"
	"github.com/verifa/coastline/requests"
)

const (
	subjectTriggerRun   = "trigger.run"
	subjectTaskResponse = "task.response"
	queueWorker         = "worker"
)

type Config struct{}

func Start(engine *requests.Engine, config *Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	worker, err := newWorker(engine)
	if err != nil {
		return fmt.Errorf("creating worker: %w", err)
	}

	if err := worker.subscribe(); err != nil {
		return fmt.Errorf("subscribing worker to nats: %w", err)
	}

	return nil
}

type worker struct {
	engine *requests.Engine
	nc     *nats.Conn
}

func newWorker(engine *requests.Engine) (*worker, error) {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to nats: %w", err)
	}
	return &worker{
		engine: engine,
		nc:     nc,
	}, nil
}

func (w *worker) subscribe() error {
	_, err := w.nc.QueueSubscribe(subjectTriggerRun, queueWorker, func(msg *nats.Msg) {
		var triggerMsg TriggerMsg
		if err := json.Unmarshal(msg.Data, &triggerMsg); err != nil {
			log.Printf("Error: unmarshalling message: %s", err.Error())
			return
		}

		if err := w.handleTrigger(&triggerMsg); err != nil {
			log.Printf("Error: handling trigger: %s", err.Error())
		}
	})
	return err
}

func (w *worker) handleTrigger(msg *TriggerMsg) error {

	resp := ResponseMsg{
		TriggerID: msg.TriggerID,
	}
	if msg.Request == nil {
		resp.Error = "request is nil"
		return w.publishResponse(&resp)
	}
	tasks, err := w.engine.GetTasksForRequest(msg.Request)
	if err != nil {
		resp.Error = fmt.Sprintf("getting tasks for request: %s", err.Error())
		return w.publishResponse(&resp)
	}
	if len(tasks) == 0 {
		resp.Error = "no tasks found"
		return w.publishResponse(&resp)
	}

	runner := requests.NewRunner(msg.Request)

	// Loop through tasks and run them, publishing results for each.
	// TODO: should we catch errors from publishing response and return those as at the end?
	for _, task := range tasks {
		taskName, ok := task.Label()
		if !ok {
			continue
		}
		resp.Task = taskName
		output, err := runner.RunTask(task)
		if err != nil {
			resp.Error = fmt.Sprintf("running task: %s", err.Error())
			w.publishResponse(&resp)
			continue
		}

		if err := output.Validate(cue.Concrete(true)); err != nil {
			resp.Error = fmt.Sprintf("Output invalid: %s", err.Error())
			w.publishResponse(&resp)
			continue
		}
		resp.Output, err = output.MarshalJSON()
		if err != nil {
			resp.Error = fmt.Sprintf("marshalling json: %s", err.Error())
			w.publishResponse(&resp)
			continue
		}

		w.publishResponse(&resp)
	}

	return nil
}

func (w *worker) publishResponse(resp *ResponseMsg) error {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error: marshalling response: %s", err.Error())

		// Create a new message to marshal
		respBytes, _ = json.Marshal(&ResponseMsg{
			TriggerID: resp.TriggerID,
			Error:     fmt.Sprintf("Error: marshalling response: %s", err.Error()),
		})
	}
	return w.nc.Publish(subjectTaskResponse, respBytes)
}
