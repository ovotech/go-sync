//nolint:gochecknoglobals,gochecknoinits
package main

import (
	"github.com/spf13/cobra"

	"github.com/ovotech/go-sync/cmd/plan"
)

var rootCmd = &cobra.Command{
	Use:   "gosync",
	Short: "Go Sync all the things",
	Long: `You have people*. You have places you want them to be. Go Sync makes it happen.
* Doesn't have to be people.

More information at https://github.com/ovotech/go-sync`,
}

func init() {
	rootCmd.AddCommand(plan.Plan)
}

func main() {
	_ = rootCmd.Execute()
}
