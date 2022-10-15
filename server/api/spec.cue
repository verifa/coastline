package api

#Project: {
	name: =~"[A-Za-z0-9-].*"
}

#Service: {
	name: string
}

#Request: {
	project:   #Project
	service:   #Service
	requester: string
}

service: #Service
