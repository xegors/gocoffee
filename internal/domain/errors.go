package domain

import "errors"

var (
	// ErrUnknownCategory is returned by ParseCategory when the value matches
	// none of the known categories. Callers match it with errors.Is.
	ErrUnknownCategory = errors.New("unknown category")

	// ErrUnknownBrew is returned by ParseBrewMethod when the value matches
	// none of the known roast methods. An empty string is not an error: it
	// resolves to BrewNone.
	ErrUnknownBrew = errors.New("unknown brew method")
)
