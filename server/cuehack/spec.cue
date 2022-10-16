package cuehack

#Requests: #ArtifactoryRepoRequest | #JenkinsServerRequest

#Request: {
	project_id:   string
	service_id:   string
	requested_by: string
}

#ArtifactoryRepoRequest: {
	#Request
	type: "ArtifactoryRepoRequest"
	spec: {
		repo: string
	}
}

#JenkinsServerRequest: {
	#Request
	type: "JenkinsServerRequest"
	spec: {
		name: string
	}
}

request: #Requests
