package basic

import (
	"encoding/json"
	"github.com/verifa/coastline/tasks/http"
)

#CatFact: {
	kind: "CatFact"
	service: {
		selector: {
			matchLabels: {
				tool: "catfact"
			}
		}
	}
	spec: {
		// Max length of cat fact
		maxLength: int | *100
	}
}

workflow: CatFact: {
	input: #CatFact

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
