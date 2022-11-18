package demo

test: PokemonDitto: {
	run: workflow.PokemonInfo & {
		input: spec: {
			name: "ditto"
		}
	}

	assert: nationalDexID: run.output.national_dex_id == 132
}
