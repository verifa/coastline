package demo

import (
	"encoding/json"
	
	"github.com/verifa/coastline/tasks/http"
)

request: #PokemonInfo: {
	kind:        "PokemonInfo"
	description: "Info on \(spec.name)"
	serviceSelector: {
		matchLabels: {
			tool: "pokemon"
		}
	}
	spec: {
		name: string
	}
}

workflow: PokemonInfo: {
	input: request.#PokemonInfo

	step: api: http.Get & {
		url: "https://pokeapi.co/api/v2/pokemon/\(input.spec.name)"
	}

	step: answer: json.Unmarshal(step.api.response.body)

	output: {
		name:            step.answer.name
		national_dex_id: step.answer.id
		height:          step.answer.height
		types:           step.answer.types
		weight:          step.answer.weight
	}
}
