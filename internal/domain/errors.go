package domain

import "errors"

var (
	// ErrUnknownCategory is returned by ParseCategory when the value matches
	// none of the known categories. Callers match it with errors.Is.
	ErrUnknownCategory = errors.New("unknown category")

	// ErrInvalidAmount is returned by Deposit or Withdraw when the amount is
	// zero or negative.
	ErrInvalidAmount = errors.New("invalid amount")

	// ErrInsufficientFunds is returned by Withdraw when the requested amount
	// exceeds the wallet balance.
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrInvalidItem is returned by Cart.Add when the item has no ID.
	ErrInvalidItem = errors.New("invalid item")

	// ErrInvalidQuantity is returned by Cart.Add or Cart.Remove when the
	// requested quantity is non-positive, or when adding it would overflow int.
	ErrInvalidQuantity = errors.New("invalid quantity")

	// ErrItemNotInCart is returned by Cart.Remove when the requested item is
	// not present in the cart.
	ErrItemNotInCart = errors.New("item is not in cart")

	// ErrInsufficientQuantity is returned by Cart.Remove when the requested
	// quantity exceeds the quantity of that item in the cart.
	ErrInsufficientQuantity = errors.New("insufficient quantity")

	// ErrInvalidUserID is returned by User.Validate when a user carries an
	// empty ID. Such a user cannot be looked up later, so storing it is
	// refused instead of accepting a record nobody can reach.
	ErrInvalidUserID = errors.New("invalid user id")
)
