package domain

import "math"

// Money is a whole number of rubles.
//
// Zero kopecks are representable: 180 means 180 rubles, and prices in this
// shop are set in whole rubles. The type is a ruble counter, not a
// multi-currency amount — a crypto payer converts into its own units and
// its own type, and never into Money.
//
// The zero value is zero rubles.
type Money int64

// addMoney returns the sum of a and b.
//
// It reports false when either operand is negative, or when the sum does not
// fit in Money.
func addMoney(a, b Money) (Money, bool) {
	switch {
	case a < 0 || b < 0:
		return 0, false
	case b > math.MaxInt64-a:
		return 0, false
	default:
		return a + b, true
	}
}
