package main

import (
	"github.com/ovotech/go-sync/adapters/test/inmem"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(inmem.Init)
}
