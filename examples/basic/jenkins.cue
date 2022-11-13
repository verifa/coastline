package basic

#JenkinsServerRequest: {
	kind: "JenkinsServerRequest"
	service: {
		selector: {
			matchLabels: {
				tool: "artifactory"
			}
		}
	}
	spec: {
		// Name must not contain space or strange characters
		name: =~"^[A-Za-z0-9-]+$"
	}
}

workflow: jenkinsServer: {
	input: #JenkinsServerRequest

	output: {
		server: "awesome-jenkins-server-\(input.spec.name)"
	}
}
