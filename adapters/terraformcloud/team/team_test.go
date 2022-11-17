package team_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
)

func ExampleNew() {
	adapter := team.New()

	gosync.New(adapter)
}

func ExampleInit() {
	ctx := context.Background()

	adapter, err := team.Init(ctx, map[gosync.ConfigKey]string{team.AnExampleConfig: "example"})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
