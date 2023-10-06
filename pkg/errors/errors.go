package errors

import "errors"

// ErrNotImplemented is for brand-new adapters that are still being worked on.
//
//goland:noinspection GoUnusedGlobalVariable
var ErrNotImplemented = errors.New("not implemented")

// ErrCacheEmpty is returned if an adapter expects Get to be called before Add/Remove.
var ErrCacheEmpty = errors.New("cache is empty, run Get first")

// ErrReadOnly is returned for adapters that cannot Add/Remove, but have been set as a destination.
var ErrReadOnly = errors.New("cannot perform action, adapter is readonly")

// ErrMissingConfig is returned when an InitFn is missing a required configuration.
var ErrMissingConfig = errors.New("missing configuration")

// ErrInvalidConfig is returned when an InitFn is passed an invalid configuration.
var ErrInvalidConfig = errors.New("invalid configuration")

// ErrTooManyChanges is returned when a change limit has been set, and the number of changes exceeds it.
var ErrTooManyChanges = errors.New("too many changes")

// ErrDoesNotExist is returned when something does not exist.
var ErrDoesNotExist = errors.New("does not exist")
