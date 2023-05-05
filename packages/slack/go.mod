module github.com/ovotech/go-sync/packages/slack

go 1.18

require (
	github.com/ovotech/go-sync/packages/gosync v0.0.0
	github.com/slack-go/slack v0.12.2
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ovotech/go-sync/packages/gosync v0.0.0 => ../gosync
