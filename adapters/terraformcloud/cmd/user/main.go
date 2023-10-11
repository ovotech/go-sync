package main

import (
	"github.com/ovotech/go-sync/adapters/terraformcloud/user"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(user.Init)
}
