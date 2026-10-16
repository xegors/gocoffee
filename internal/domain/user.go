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
// Deposit and Withdraw cannot make the balance negative: both reject an amount
// that would cross zero. A struct literal can still set a negative balance,
// since the field is exported, and Deposit refuses to grow such a balance.
// Failed deposits and withdrawals leave the balance unchanged.
type Wallet struct {
	Balance Money
}

// Deposit adds amount to the wallet balance.
//
// It returns ErrInvalidAmount and leaves the balance unchanged when amount is
// zero or negative, when the balance is negative already, or when the sum does
// not fit in Money.
//
// The balance cannot become negative through this method, but it can be set
// to a negative value through a struct literal: the check on both operands
// keeps such a balance from being grown further.
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

// User represents a customer of the shop, together with everything the shop
// remembers about them: who they are, what they can spend, and how much they
// have already spent.
//
// The wallet is embedded by value so its fields and methods are promoted,
// allowing access through user.Balance and user.Deposit(...).
//
// A user carries no cart. A cart is a separate piece of session state, and
// Stage 1 has nowhere to keep it: Cart holds its lines in an unexported slice
// that cannot even be read from outside the package.
type User struct {
	ID string
	Wallet
	// Email is the address the customer orders under.
	//
	// Nothing checks its format: the shop has no rule for addresses yet, and
	// the rule arrives with authentication on Stage 2. An address that is
	// stored but never validated only accumulates typos. The empty string is a
	// valid value — a customer created without an address is still a customer.
	Email string
	// RegisteredAt is when the customer first appeared.
	//
	// The zero value means "date unknown" and is a valid value: files written
	// before this field existed, and records filled in by hand, load without
	// failing the whole file.
	RegisteredAt time.Time
	// TotalSpent is the sum of every payment the customer has ever made.
	//
	// It only grows: a discount lowers the price of an order, never the money
	// already spent. The loyalty rank is derived from this value with
	// RankFor and is deliberately not stored — a stored rank would be a
	// second copy of the same fact, free to disagree with the sum.
	//
	// The zero value means the customer has never paid.
	TotalSpent Money
}

// Validate reports whether the user can be stored.
//
// It checks the ID and nothing else: an empty ID, or one made only of
// whitespace, is ErrInvalidUserID. Email, dates and amounts are stored as
// given, because the shop has no rules about them yet — a rule that is missing
// is visible, a rule that is checked loosely hides the data it was meant to
// keep out.
//
// A whitespace-only ID is rejected on purpose: it is not a name anybody can
// type, and it is indistinguishable from "forgot to fill it in".
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
// not fit in Money. The reasoning is the same as in Wallet.Deposit: a failed
// accounting entry must not shift the number the discount is computed from.
//
// There is no way to lower the total. Returns and refunds are not part of
// Stage 1, and a discount reduces the price of an order rather than the
// history of payments — otherwise a customer could farm the loyalty rank by
// buying and unwinding.
func (u *User) AddSpent(m Money) error {
	if m > 0 {
		if sum, ok := addMoney(u.TotalSpent, m); ok {
			u.TotalSpent = sum
			return nil
		}
	}

	return ErrInvalidAmount
}
