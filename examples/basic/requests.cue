package basic

#Request: #ArtifactoryRepoRequest | #JenkinsServerRequest

#metadata: {
	name:      string
	namespace: string
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
