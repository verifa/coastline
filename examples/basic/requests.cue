package basic

#Request: #ArtifactoryRepoRequest | #JenkinsServerRequest
request:  #Request

#BaseRequest: {
	project_id:   string
	service_id:   string
	requested_by: string
}

#ArtifactoryRepoRequest: {
	#BaseRequest
	type: "ArtifactoryRepoRequest"
	spec: {
		repo: string
	}
}

#JenkinsServerRequest: {
	#BaseRequest
	type: "JenkinsServerRequest"
	spec: {
		name: string
	}
}
