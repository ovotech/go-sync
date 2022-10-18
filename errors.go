package gosync

import "errors"

// ErrNotImplemented is for brand-new adapters that are still being worked on.
var ErrNotImplemented = errors.New("not implemented")

// ErrCacheEmpty is returned if an adapter expects Get to be called before Add/Remove.
var ErrCacheEmpty = errors.New("cache is empty, run Get first")

// ErrReadOnly is returned for adapters that cannot Add/Remove, but have been set as a destination.
var ErrReadOnly = errors.New("cannot perform action, adapter is readonly")

// ErrMissingConfig is returned when an InitFn is missing a required configuration.
var ErrMissingConfig = errors.New("missing configuration")

// ErrInvalidConfig is returned when an InitFn is passed an invalid configuration.
var ErrInvalidConfig = errors.New("invalid configuration")
