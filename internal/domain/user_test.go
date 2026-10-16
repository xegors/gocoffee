package domain_test

import (
	"errors"
	"math"
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
		{name: "negative balance is not grown", initial: -100, amount: 50,
			wantBalance: -100, wantErr: domain.ErrInvalidAmount},
		{name: "sum at the money limit", initial: math.MaxInt64 - 1, amount: 1,
			wantBalance: math.MaxInt64},
		{name: "sum over the money limit", initial: math.MaxInt64, amount: 1,
			wantBalance: math.MaxInt64, wantErr: domain.ErrInvalidAmount},
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

func TestWalletPay(t *testing.T) {
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
			err := w.Pay(tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("|got %10q| want %10q|", err, tt.wantErr)
			}
			if w.Balance != tt.wantBalance {
				t.Errorf("|got %10d| want %10d|", w.Balance, tt.wantBalance)
			}
		})
	}
}

func TestAddSpent(t *testing.T) {
	tests := []struct {
		name           string
		initial        domain.Money
		amount         domain.Money
		balance        domain.Money
		wantTotalSpent domain.Money
		wantErr        error
	}{
		{
			name: "zero amount", initial: 1000, amount: 0,
			wantTotalSpent: 1000, wantErr: domain.ErrInvalidAmount,
		},
		{
			name: "negative amount", initial: 1000, amount: -200,
			wantTotalSpent: 1000, wantErr: domain.ErrInvalidAmount,
		},
		{
			name: "negative total is not grown", initial: -100, amount: 50,
			wantTotalSpent: -100, wantErr: domain.ErrInvalidAmount,
		},
		{
			// Кошелёк и потраченное — разные суммы: ошибка в выборе
			// операнда сложила бы их вместе и потеряла историю покупок.
			name: "wallet balance is not mixed in", initial: 1000, amount: 200, balance: 5555,
			wantTotalSpent: 1200,
		},
		{
			name: "MaxInt64 border", initial: math.MaxInt64 - 1, amount: 1,
			wantTotalSpent: math.MaxInt64,
		},
		{
			name: "over MaxInt64", initial: math.MaxInt64, amount: 1,
			wantTotalSpent: math.MaxInt64, wantErr: domain.ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := domain.User{TotalSpent: tt.initial, Wallet: domain.Wallet{Balance: tt.balance}}
			if err := u.AddSpent(tt.amount); !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}

			if u.TotalSpent != tt.wantTotalSpent {
				t.Errorf("got %d, want %d", u.TotalSpent, tt.wantTotalSpent)
			}
		})
	}
}

func TestUserValidate(t *testing.T) {
	tests := []struct {
		name    string
		ID      string
		wantErr error
	}{
		{name: "empty ID", ID: "", wantErr: domain.ErrInvalidUserID},
		{name: "regular ID", ID: "some_normal_id"},
		{name: "only spaces ID", ID: " \t \t ", wantErr: domain.ErrInvalidUserID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := domain.User{ID: tt.ID}
			if err := u.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

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
