package main

import (
	"github.com/ovotech/go-sync/adapters/azuread/groupmembership"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(groupmembership.Init)
}
