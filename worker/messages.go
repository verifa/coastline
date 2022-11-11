package worker

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/verifa/coastline/server/oapi"
)

type TriggerMsg struct {
	TriggerID uuid.UUID     `json:"trigger_id"`
	Request   *oapi.Request `json:"request"`
}

type ResponseMsg struct {
	TriggerID uuid.UUID       `json:"trigger_id"`
	Workflow  string          `json:"workflow"`
	Error     string          `json:"error,omitempty"`
	Output    json.RawMessage `json:"output"`
}
