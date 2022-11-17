package http

import (
	"fmt"
	"io"
	"net/http"

	"cuelang.org/go/tools/flow"
	"github.com/verifa/coastline/tasks/helper"
)

var _ flow.Runner = (*DoTask)(nil)

type DoTask struct{}

func (t DoTask) Run(task *flow.Task, pErr error) error {
	th := helper.TaskHelper{
		Task: task,
	}

	method := th.MustString("method")
	url := th.MustString("url")
	body, _ := th.Reader("request.body")

	var params map[string]string
	th.Decode("request.params", &params)

	// Check if there were any errors with config before proceeding
	if th.Errs != nil {
		return th.Errs
	}

	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
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
