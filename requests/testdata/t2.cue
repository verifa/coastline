package testdata

#T2: {
	kind: "t2"
	service: {
		selector: {
			matchLabels: {
				key: "value"
			}
		}
	}
	spec: {
		// Name must not contain space or strange characters
		num: number
	}
}

workflow: t2: {
	input: #T2

	output: {
		key: input.spec.num
	}
}
