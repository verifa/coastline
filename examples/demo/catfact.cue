package demo

import (
	"encoding/json"
	
	"github.com/verifa/coastline/tasks/http"
)

request: #CatFact: {
	kind:        "CatFact"
	description: "Cat fact max length \(spec.maxLength)"
	serviceSelector: {
		matchLabels: {
			tool: "cat-facts"
		}
	}
	spec: {
		// Max length of cat fact
		maxLength: int | *100
	}
}

workflow: CatFact: {
	input: request.#CatFact

	step: api: http.Get & {
		url: "https://catfact.ninja/fact"
		request: {
			params: {
				max_length: "\(input.spec.maxLength)"
			}
		}
	}

	output: {
		fact: json.Unmarshal(step.api.response.body).fact
	}
}
