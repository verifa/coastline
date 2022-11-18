package basic

import "tool/http"

request: #ExampleRequest: {
	kind: "ExampleRequest"
	spec: {
		name: =~"^[A-Za-z0-9-]+$"
	}
}

workflow: example: {
	input: request.#ExampleRequest

	step: restAPI: http.Get & {
		url: "<url>"
	}

	output: {
		status: step.restAPI.response.status
	}
}
