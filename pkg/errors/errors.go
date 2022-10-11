/*
Package errors contains errors returned from Go Sync and adapters.
*/
package errors

import "errors"

// ErrCacheEmpty is returned if an adapter expects Get to be called before Add/Remove.
var ErrCacheEmpty = errors.New("cache is empty, run Get first")

// ErrReadOnly is returned for adapters that cannot Add/Remove, but have been set as a destination.
var ErrReadOnly = errors.New("cannot perform action, adapter is readonly")
