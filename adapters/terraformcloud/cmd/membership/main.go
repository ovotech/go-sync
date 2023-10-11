package main

import (
	"github.com/ovotech/go-sync/adapters/terraformcloud/membership"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(membership.Init)
}
