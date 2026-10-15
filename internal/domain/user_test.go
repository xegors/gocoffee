package domain_test

import (
	"errors"
	"testing"

	"github.com/xegors/gocoffee/internal/domain"
)

func TestWalletDeposit(t *testing.T) {
	tests := []struct {
		name        string
		initial     domain.Money
		amount      domain.Money
		wantBalance domain.Money
		wantErr     error
	}{
		{name: "to zero deposit", initial: 0, amount: 20, wantBalance: 20},
		{name: "reg deposit", initial: 2342, amount: 1231, wantBalance: 2342 + 1231},
		{name: "negative deposit", initial: 1000, amount: -500, wantBalance: 1000,
			wantErr: domain.ErrInvalidAmount},
		{name: "zero deposit", initial: 5000, amount: 0, wantBalance: 5000,
			wantErr: domain.ErrInvalidAmount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := domain.Wallet{Balance: tt.initial}
			err := w.Deposit(tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("|got %10q| want %10q|", err, tt.wantErr)
			}
			if w.Balance != tt.wantBalance {
				t.Errorf("|got %10d| want %10d|", w.Balance, tt.wantBalance)
			}
		})
	}
}

func TestWalletWithdraw(t *testing.T) {
	tests := []struct {
		name        string
		initial     domain.Money
		amount      domain.Money
		wantBalance domain.Money
		wantErr     error
	}{
		{name: "regular withdraw", initial: 4000, amount: 900, wantBalance: 4000 - 900},
		{name: "withdraw all balance", initial: 1000, amount: 1000, wantBalance: 0},
		{name: "withdraw over a balance", initial: 1000, amount: 1500, wantBalance: 1000,
			wantErr: domain.ErrInsufficientFunds},
		{name: "zero withdraw", initial: 1000, amount: 0, wantBalance: 1000,
			wantErr: domain.ErrInvalidAmount},
		{name: "negative withdraw", initial: 2340, amount: -200, wantBalance: 2340,
			wantErr: domain.ErrInvalidAmount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := domain.Wallet{Balance: tt.initial}
			err := w.Withdraw(tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("|got %10q| want %10q|", err, tt.wantErr)
			}
			if w.Balance != tt.wantBalance {
				t.Errorf("|got %10d| want %10d|", w.Balance, tt.wantBalance)
			}
		})
	}
}

// TestUserPromotesWallet:
// - поле Balance доступно как user.Balance;
// - Deposit и Withdraw вызываются через User;
// - нулевой Wallet внутри User безопасен.
func TestUserPromotesWallet(t *testing.T) {
	u := domain.User{ID: "bob"}
	err := u.Deposit(100)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}
	if u.Balance != 100 {
		t.Errorf("|got %10d| want %10d|", u.Balance, 100)
	}

	err = u.Withdraw(50)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}
	if u.Balance != 50 {
		t.Errorf("|got %10d| want %10d|", u.Balance, 50)
	}

	u = domain.User{}
	if err = u.Withdraw(0); !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("|got %10q| want %10q|", err, domain.ErrInvalidAmount)
	}
}
