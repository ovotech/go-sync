//nolint:gochecknoglobals
package plan

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ovotech/go-sync/internal/config"
)

var Plan = &cobra.Command{
	Use:   "plan",
	Short: "Dry run a Go Sync configuration",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(args[0])
		if err != nil {
			cmd.PrintErr(err)

			return
		}

		cfg.Run(context.Background(), true)
	},
}
