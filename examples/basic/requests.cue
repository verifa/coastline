package basic

#Request: #ArtifactoryRepoRequest | #JenkinsServerRequest

#metadata: {
	name:      string
	namespace: string
}

// An example request with simple types 1
#SimpleTypes: {
	// An example request with simple types 2
	type: "SimpleTypes"
	spec: {
		text:          string
		num:           number
		integer:       uint8
		boolean:       bool
		defaultString: string | *"default"
		stringArray: [...string]
		numberArray: [...number]
		intArray: [...uint8]
		nestedArray: [...[...string]]
		nested: {
			nestedText:        string
			stringEnum:        "yes" | "no" | "maybe" | "perhaps"
			numberEnum:        0 | 50 | 100
			numberEnumDefault: 0 | *50 | 100
		}
	}
}

#ArtifactoryRepoRequest: {
	type:    "ArtifactoryRepoRequest"
	service: "artifactory"
	spec: {
		repo:     string
		metadata: #metadata
	}
}

#JenkinsServerRequest: {
	type: "JenkinsServerRequest"
	spec: {
		name: string
	}
}
