package main

import (
	"github.com/ovotech/go-sync/adapters/terraformcloud/team"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(team.Init)
}
