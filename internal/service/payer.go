package service

import "github.com/xegors/gocoffee/internal/domain"

// Payer is a payment method that can charge a customer an amount.
//
// Declared in the service package because Checkout, its only consumer,
// is what needs to know about payment methods. The domain does not import
// this package: Wallet satisfies Payer structurally, without referencing it.
type Payer interface {
	Pay(amount domain.Money) error
}
