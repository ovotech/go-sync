module github.com/ovotech/go-sync/packages/github

go 1.18

require (
	github.com/google/go-github/v47 v47.1.0
	github.com/ovotech/go-sync/packages/gosync v0.0.0
	github.com/shurcooL/githubv4 v0.0.0-20220520033151-0b4e3294ff00
	github.com/stretchr/testify v1.8.4
	golang.org/x/oauth2 v0.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shurcooL/graphql v0.0.0-20220606043923-3cf50f8a0a29 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ovotech/go-sync/packages/gosync v0.0.0 => ../gosync
