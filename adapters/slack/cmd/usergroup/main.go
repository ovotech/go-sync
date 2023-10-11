package main

import (
	"github.com/ovotech/go-sync/adapters/slack/usergroup"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(usergroup.Init)
}
