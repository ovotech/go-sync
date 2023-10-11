package main

import (
	"github.com/ovotech/go-sync/adapters/github/team"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(team.Init)
}
