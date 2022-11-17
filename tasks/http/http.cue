package http

Get:    Do & {method: "GET"}
Post:   Do & {method: "POST"}
Put:    Do & {method: "PUT"}
Delete: Do & {method: "DELETE"}

Do: {
	$task: "http.Do"

	method: string
	url:    string

	request: {
		body?: bytes | string
		header: [string]:  string | [...string]
		trailer: [string]: string | [...string]
	}
	response: {
		status:     string
		statusCode: int
		body:       bytes
		header: [string]:  string | [...string]
		trailer: [string]: string | [...string]
	}
}
