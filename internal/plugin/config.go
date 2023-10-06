package plugin

import (
	goplugin "github.com/hashicorp/go-plugin"
)

//nolint:gochecknoglobals
var HandshakeConfig = goplugin.HandshakeConfig{
	ProtocolVersion:  0,
	MagicCookieKey:   "Service",
	MagicCookieValue: "GoSync",
}
