package demo

test: CatFact: {
	run: workflow.CatFact & {
		input: spec: {
			maxLength: 50
		}
	}

	assert: statusCode: run.step.api.response.statusCode == 200
	assert: factLength: len(run.output.fact) <= 50
}
