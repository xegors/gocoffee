package storage

import "errors"

var (
	// ErrCorruptedCatalog is returned when the catalog file exists and can be
	// read, but its content is not a valid catalog. It is wrapped with the file
	// path and the underlying json error, so callers should use errors.Is to
	// detect it.
	ErrCorruptedCatalog = errors.New("corrupted catalog file")
)
