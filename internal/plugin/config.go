//nolint:gochecknoglobals
package plugin

import (
	goplugin "github.com/hashicorp/go-plugin"
)

const AdapterName = "adapter"

var HandshakeConfig = goplugin.HandshakeConfig{
	ProtocolVersion:  0,
	MagicCookieKey:   "Service",
	MagicCookieValue: "GoSync",
}
