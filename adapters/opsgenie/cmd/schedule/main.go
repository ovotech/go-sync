package main

import (
	"github.com/ovotech/go-sync/adapters/opsgenie/schedule"
	"github.com/ovotech/go-sync/pkg/plugin"
)

func main() {
	plugin.Serve(schedule.Init)
}
