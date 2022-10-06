package types

import "time"

// Logger is a subset of log.Logger to allow compatible loggers to be used with Sync.
type Logger interface {
	Println(v ...any)
	Printf(format string, v ...any)
}

// Clock is a subset of time.Time which allows us to mock the clock in tests.
type Clock interface {
	Now() time.Time
}
