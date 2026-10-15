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

	// ErrInvalidAmount is returned by Deposit or Withdraw when the amount is
	// zero or negative.
	ErrInvalidAmount = errors.New("invalid amount")

	// ErrInsufficientFunds is returned by Withdraw when the requested amount
	// exceeds the wallet balance.
	ErrInsufficientFunds = errors.New("insufficient funds")
)
