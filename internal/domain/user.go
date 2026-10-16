package domain

import (
	"strings"
	"time"
)

// Wallet holds the money a user can spend.
//
// It has one balance and no bonus balance: without separate spending rules,
// a second balance would only duplicate the same concept.
//
// The balance cannot go negative. Failed deposits and withdrawals leave it
// unchanged.
type Wallet struct {
	Balance Money
}

// Deposit adds amount to the wallet balance.
//
// It returns ErrInvalidAmount when amount is zero or negative, when the
// balance is negative already, or when the sum does not fit in Money.
func (w *Wallet) Deposit(m Money) error {
	if m > 0 {
		if sum, ok := addMoney(w.Balance, m); ok {
			w.Balance = sum
			return nil
		}
	}

	return ErrInvalidAmount
}

// Withdraw removes amount from the wallet balance.
//
// It returns ErrInvalidAmount when amount is zero or negative, and
// ErrInsufficientFunds when amount exceeds the balance. The balance remains
// unchanged if either error is returned.
func (w *Wallet) Withdraw(m Money) error {
	if m <= w.Balance && m > 0 {
		w.Balance -= m
		return nil
	}
	if m <= 0 {
		return ErrInvalidAmount
	}
	return ErrInsufficientFunds
}

// Pay charges amount from the wallet.
//
// It delegates to Withdraw, so manual withdrawal and wallet payment cannot diverge.
func (w *Wallet) Pay(m Money) error { return w.Withdraw(m) }

// User represents a customer of the shop.
//
// The wallet is embedded by value so its fields and methods are promoted,
// allowing access through user.Balance and user.Deposit(...).
type User struct {
	ID string
	Wallet
	// Email is the address the customer orders under. It is a free string:
	// the shop has no rule for addresses yet, and the rule arrives with
	// authentication on Stage 2. The empty string is a valid value.
	Email string
	// RegisteredAt is when the customer first appeared. The zero value means
	// "date unknown" and is valid: files written before the field existed load
	// without failing.
	RegisteredAt time.Time
	// TotalSpent is the sum of every payment the customer has ever made. It
	// only grows, and the loyalty rank is derived from it with RankFor rather
	// than stored next to it. Zero means the customer has never paid.
	TotalSpent Money
}

// Validate reports whether the user can be stored.
//
// It checks the ID and nothing else: an empty ID, or one made only of
// whitespace, is ErrInvalidUserID. Email, dates and amounts are stored as
// given, because the shop has no rules about them yet.
func (u User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return ErrInvalidUserID
	}

	return nil
}

// AddSpent adds amount to the money the customer has spent in total.
//
// It returns ErrInvalidAmount and leaves TotalSpent unchanged when amount is
// zero or negative, when TotalSpent is negative already, or when the sum does
// not fit in Money.
//
// There is no way to lower the total: a discount lowers the price of an order,
// never the history of payments.
func (u *User) AddSpent(m Money) error {
	if m > 0 {
		if sum, ok := addMoney(u.TotalSpent, m); ok {
			u.TotalSpent = sum
			return nil
		}
	}

	return ErrInvalidAmount
}
