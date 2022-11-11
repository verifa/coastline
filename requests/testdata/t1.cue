package testdata

import (
	"encoding/json"
	
	"github.com/verifa/coastline/tasks/http"
)

#T1: {
	kind: "t1"
	service: {
		selector: {
			matchLabels: {
				tool: "t1"
			}
		}
	}
	spec: {
		// Name must not contain space or strange characters
		name: =~"^[A-Za-z0-9-]+$"
	}
}

workflow: t1: {
	input: #T1

	step: api: http.Get & {
		url: "https://catfact.ninja/fact"
	}

	output: json.Unmarshal(step.api.response.body).fact
}
