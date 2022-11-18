package testdata

request: #t2: {
	kind: "t2"
	serviceSelector: {
		matchLabels: {
			key: "value"
		}
	}
	spec: {
		// Name must not contain space or strange characters
		num: number
	}
}

workflow: t2: {
	input: request.#t2

	output: {
		key: input.spec.num
	}
}
