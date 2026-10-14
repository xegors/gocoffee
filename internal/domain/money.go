package domain

// Money is a whole number of rubles.
//
// Zero kopecks are representable: 180 means 180 rubles, and prices in this
// shop are set in whole rubles. The type is a ruble counter, not a
// multi-currency amount — a crypto payer converts into its own units and
// its own type, and never into Money.
//
// The zero value is zero rubles.
type Money int64
