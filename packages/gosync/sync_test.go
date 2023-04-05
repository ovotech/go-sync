package gosync_test

import (
	"context"
	"log"

	"github.com/ovotech/go-sync/packages/gosync"
)

func ExampleNew() {
	// Any Go Sync adapter.
	var source gosync.Adapter

	gosync.New(source)
}

func ExampleSetCaseSensitive() {
	// Any Go Sync adapter.
	var source gosync.Adapter

	gosync.New(source, gosync.SetCaseSensitive(true))
}

func ExampleSetMaximumChanges() {
	// Any Go Sync adapter.
	var source gosync.Adapter

	gosync.New(source, gosync.SetMaximumChanges(5))
}

func ExampleSetOperatingMode() {
	// Any Go Sync adapter.
	var source gosync.Adapter

	// Set the operating mode to add only (don't remove things).
	operatingMode := gosync.SetOperatingMode(gosync.AddOnly)

	gosync.New(source, operatingMode)
}

func ExampleSync_SyncWith() {
	// Any Go Sync packages.
	var source, destination gosync.Adapter

	sync := gosync.New(source)

	// By default, Go Sync runs in dry run mode. To make changes this must manually be set to false.
	sync.DryRun = false

	err := sync.SyncWith(context.Background(), destination)
	if err != nil {
		log.Panic(err)
	}
}
