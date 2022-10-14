/*
This is an example application to demonstrate querying the user info endpoint.
*/
package main

import (
	"context"
	"go-authnz-server/server"
	"log"
	"net/http"
)

func main() {
	ctx := context.TODO()
	s, err := server.New(ctx)
	if err != nil {
		log.Fatalf("creating server: %s", err.Error())
	}

	http.ListenAndServe(":3000", s)
}
