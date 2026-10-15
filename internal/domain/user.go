package domain

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

// User represents a customer of the shop.
//
// The wallet is embedded by value so its fields and methods are promoted,
// allowing access through user.Balance and user.Deposit(...).
type User struct {
	ID string
	Wallet
}

// Deposit adds amount to the wallet balance.
//
// It returns ErrInvalidAmount when amount is zero or negative and leaves the
// balance unchanged.
func (w *Wallet) Deposit(m Money) error {
	if m > 0 {
		w.Balance += m
		return nil
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
