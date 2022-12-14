package basic

#metadata: {
	name:      string
	namespace: string
}

// An example request with simple types 1
request: #SimpleTypes: {
	kind: "SimpleTypes"
	serviceSelector: {
		matchLabels: {
			tool: "artifactory"
		}
	}
	spec: {
		// Just a text field
		text: string
		// Simply a number field
		num: number
		// Provide an integer, not a number
		integer: uint8
		// Boolean, true or false
		boolean: bool
		// A string enum
		stringEnum: "yes" | "no" | "maybe" | "perhaps"
		// String with a default
		defaultString: string | *"default"
		// Array of strings
		stringArray: [...string]
		// Array of some numbers (not integers)
		numberArray: [...number]
		// Array of integers
		intArray: [...uint8]
		// An array of arrays of string
		nestedArray: [...[...string]]
		// A nested object within the example
		nested: {
			nestedText:        string
			stringEnum:        "yes" | "no" | "maybe" | "perhaps"
			numberEnum:        0 | 50 | 100
			numberEnumDefault: 0 | *50 | 100
		}
		anotherNested: {
			empty:   string
			default: string | *"value"
		}
	}
}

request: #ArtifactoryRepoRequest: {
	kind: "ArtifactoryRepoRequest"
	service: {
		selector: {
			matchLabels: {
				tool: "artifactory"
			}
		}
	}
	spec: {
		// Name of the repository to create
		repo: string
		// General metadata for the request
		metadata: #metadata
	}
}
