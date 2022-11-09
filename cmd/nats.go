package cmd

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

// TODO: think about location for this.
// TODO: make a struct so that shutdown can do ns.Shutdown() and ns.WaitForShutdown() on Close()

func startNats() error {
	opts := &server.Options{}
	ns, err := server.NewServer(opts)
	if err != nil {
		return err
	}

	go ns.Start()

	timeout := time.Second * 4
	// Wait for server to be ready for connections
	if !ns.ReadyForConnections(timeout) {
		return fmt.Errorf("nats not ready for connection after: %s", timeout)
	}
	return nil
}
