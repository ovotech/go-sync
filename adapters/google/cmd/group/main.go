package main

import (
	"github.com/ovotech/go-sync/adapters/google/group"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(group.Init)
}
