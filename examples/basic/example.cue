package basic

import "tool/http"

#ExampleRequest: {
	kind: "ExampleRequest"
	spec: {
		name: =~"^[A-Za-z0-9-]+$"
	}
}

task: example: {
	input: #ExampleRequest

	step: restAPI: http.Get & {
		url: "<url>"
	}

	output: {
		status: step.restAPI.response.status
	}
}
