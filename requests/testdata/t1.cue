package templates

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

task: t1: {
	input: #T1

	output: {
		key: input.spec.name
	}
}
