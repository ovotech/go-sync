package types

// Logger is a subset of log.Logger to allow compatible loggers to be used with Sync.
type Logger interface {
	Println(v ...any)
	Printf(format string, v ...any)
}
