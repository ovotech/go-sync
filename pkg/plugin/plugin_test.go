package plugin_test

import (
	"github.com/ovotech/go-sync/adapters/slack/usergroup"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func ExampleServe() {
	plugin.Serve(
		// Pass a Go Sync InitFn.
		usergroup.Init,
		// Add as many ConfigFns as you need to customise your adapter.
		usergroup.WithClient(nil),
	)
}
