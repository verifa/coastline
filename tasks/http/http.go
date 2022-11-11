package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"cuelang.org/go/tools/flow"
)

type doConfig struct {
	Method  string `json:"method"`
	URL     string `json:"url"`
	Request struct {
		Body    []byte            `json:"body"`
		Headers map[string]string `json:"headers"`
		Params  map[string]string `json:"params"`
	} `json:"request"`
}

var _ flow.Runner = (*DoTask)(nil)

type DoTask struct{}

func (t DoTask) Run(task *flow.Task, pErr error) error {
	var config doConfig
	if err := task.Value().Decode(&config); err != nil {
		return fmt.Errorf("decoding task: %w", err)
	}

	client := http.Client{}
	req, err := http.NewRequest(config.Method, config.URL, bytes.NewBuffer(config.Request.Body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if config.Request.Params != nil {
		q := req.URL.Query()
		for key, value := range config.Request.Params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	return task.Fill(map[string]interface{}{
		"response": map[string]interface{}{
			"status":     resp.Status,
			"statusCode": resp.StatusCode,
			"body":       bodyBytes,
			"header":     resp.Header,
			"trailer":    resp.Trailer,
		},
	})
}
