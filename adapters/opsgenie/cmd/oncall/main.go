package main

import (
	"github.com/ovotech/go-sync/adapters/opsgenie/oncall"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(oncall.Init)
}
