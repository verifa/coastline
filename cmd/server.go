/*
Copyright © 2022 Jacob Larfors <jlarfors@verifa.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/verifa/coastline/server"
	"github.com/verifa/coastline/store"

	"github.com/spf13/cobra"
)

var serverConfig = server.DefaultConfig()
var storeConfig store.Config

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.TODO()
		if serverConfig.DevMode {
			storeConfig.InitData = true
		}
		store, err := store.New(ctx, &storeConfig)
		if err != nil {
			return fmt.Errorf("creating store: %w", err)
		}

		srv, err := server.New(ctx, store, &serverConfig)
		if err != nil {
			return fmt.Errorf("creating server: %w", err)
		}

		return http.ListenAndServe(":3000", srv)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().BoolVarP(&serverConfig.DevMode, "dev", "d", false, "Enable dev mode")
	serverCmd.Flags().StringVarP(&serverConfig.RequestsEngine.Templates, "templates", "t", "", "Path to request templates to load")
	serverCmd.Flags().StringVarP(&serverConfig.Dir, "directory", "C", ".", "Base directory to run command from")
}
