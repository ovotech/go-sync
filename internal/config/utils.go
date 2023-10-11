package config

import (
	"os"
	"strings"

	"github.com/ovotech/go-sync/pkg/types"
)

// getEnvironmentVariables retrieves all environment variables starting with GOSYNC_.
func getEnvironmentVariables() map[types.ConfigKey]string {
	vars := os.Environ()
	out := make(map[string]string)

	for _, envVar := range vars {
		if strings.HasPrefix(envVar, "GOSYNC_") {
			if key, value, ok := strings.Cut(envVar, "="); ok {
				key, _ = strings.CutPrefix(key, "GOSYNC_")

				out[key] = value
			}
		}
	}

	return out
}
