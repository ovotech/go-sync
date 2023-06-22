package gosync_test

import (
	"context"
	"log"

	gosync "github.com/ovotech/go-sync"
)

type serviceStruct struct {
	New func(string) serviceStruct
}

type adapterStruct struct {
	Init      gosync.InitFn
	New       func(serviceStruct, string) gosync.Adapter
	Token     string
	Something string
}

//nolint:gochecknoglobals
var (
	someAdapter = adapterStruct{}
	service     = serviceStruct{}
)

func ExampleNew() {
	// Create an adapter using the recommended New method.
	client := service.New("some-token")
	source := someAdapter.New(client, "some-value")

	// Initialise an adapter using an Init function.
	destination, err := someAdapter.Init(context.Background(), map[gosync.ConfigKey]string{
		someAdapter.Token:     "some-token",
		someAdapter.Something: "some-value",
	})
	if err != nil {
		log.Fatal(err)
	}

	sync := gosync.New(source)

	err = sync.SyncWith(context.Background(), destination)
	if err != nil {
		log.Fatal(err)
	}
}
