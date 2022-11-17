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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/verifa/coastline/requests"
)

var (
	testCmdVerbose bool
	testCmdRun     string
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs CUE tests",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Enable parsing of tests
		requestsConfig.IsTesting = true
		engine, err := requests.Load(&requestsConfig)
		if err != nil {
			return fmt.Errorf("loading requests engine: %w", err)
		}
		result, err := engine.RunTests(&requests.TestConfig{
			Filter: testCmdRun,
		})
		if err != nil {
			return fmt.Errorf("running tests: %w", err)
		}
		result.Print(testCmdVerbose)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().BoolVarP(&testCmdVerbose, "verbose", "v", false, "Show output for each test")
	testCmd.Flags().StringVarP(&testCmdRun, "run", "r", "", "Filter tests to run")
}
