///usr/local/bin/dagger run go run $0 $@ ; exit

package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

var (
	client *dagger.Client
	linter *dagger.Container
	tests  *dagger.Container
)

var rootCmd = &cobra.Command{
	Use:   "./ci/dagger.go",
	Short: "A simple CLI to execute continuous-schema-migration project CI/CD with dagger",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Initialize Dagger client
		client, err = dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stderr))
		if err != nil {
			return err
		}
		defer client.Close()

		linter = Linter(client)
		tests = Tests(client)

		return nil
	},
}

var allCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a"},
	Short:   "Execute all CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), []*dagger.Container{linter, tests})
	},
}

var linterCmd = &cobra.Command{
	Use:     "linter",
	Aliases: []string{"g"},
	Short:   "Execute linter step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), []*dagger.Container{linter})
	},
}

var testsCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Execute tests step of the CI",
	Run: func(cmd *cobra.Command, args []string) {
		executeContainers(context.Background(), []*dagger.Container{tests})
	},
}

func main() {
	rootCmd.AddCommand(allCmd)
	rootCmd.AddCommand(linterCmd)
	rootCmd.AddCommand(testsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
