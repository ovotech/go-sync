module github.com/ovotech/go-sync/packages/opsgenie

go 1.18

require (
	github.com/opsgenie/opsgenie-go-sdk-v2 v1.2.14
	github.com/ovotech/go-sync/packages/gosync v0.0.0
	github.com/stretchr/testify v1.8.2
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ovotech/go-sync/packages/gosync v0.0.0 => ../gosync
