package domain

import (
	"math"
	"slices"
)

// CartItem represents one product line in a cart.
//
// Quantity is the number of units of Item in this line.
type CartItem struct {
	Item     Item
	Quantity int
}

// Cart stores product lines in insertion order.
//
// Add combines repeated additions of the same item by increasing its quantity.
// Remove subtracts the requested quantity from the matching line. Items returns
// a copy of the internal slice so callers cannot modify the cart directly.
type Cart struct {
	items []CartItem
}

// Add adds qty units of item to the cart.
//
// Repeated additions of the same item increase its existing quantity. It
// returns ErrInvalidItem for an item without an ID, and ErrInvalidQuantity
// for a non-positive quantity or an addition that would overflow int. On
// error, the cart is unchanged.
func (c *Cart) Add(item Item, qty int) error {
	if item.ID == "" {
		return ErrInvalidItem
	}
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	if i := c.indexOf(item.ID); i != -1 {
		if qty > math.MaxInt-c.items[i].Quantity {
			return ErrInvalidQuantity
		}
		c.items[i].Quantity += qty
		return nil
	}

	c.items = append(c.items, CartItem{item, qty})
	return nil
}

// Merge adds the lines in other to c.
//
// Existing item IDs have their quantities combined; new items are appended
// in the order they appear in other.
func (c *Cart) Merge(other Cart) error {
	for _, ci := range other.items {
		if err := c.Add(ci.Item, ci.Quantity); err != nil {
			return err
		}
	}
	return nil
}

// Remove subtracts qty units of the identified item from the cart.
//
// It removes the line when qty equals the current quantity. It returns
// ErrInvalidQuantity for a non-positive quantity, ErrItemNotInCart when the
// item ID is absent, and ErrInsufficientQuantity when qty exceeds the
// quantity in the cart. The cart is unchanged on error.
func (c *Cart) Remove(id string, qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	i := c.indexOf(id)
	if i == -1 {
		return ErrItemNotInCart
	}

	cur := &c.items[i]
	switch {
	case cur.Quantity == qty:
		c.items = append(c.items[:i], c.items[i+1:]...)
		return nil
	case cur.Quantity > qty:
		cur.Quantity -= qty
		return nil

	default:
		return ErrInsufficientQuantity
	}
}

// Total returns the sum of item prices multiplied by their quantities.
// Discounts are not applied here; they are handled during checkout.
//
// The sum is computed without an overflow check, unlike the money additions in
// Wallet. The two places differ deliberately: Total has no way to report a
// failure, and silently discarding one would be worse than not checking at
// all. Reaching the limit would take a single line item priced at nine
// quintillion rubles. If that ever stops being absurd, the signature changes
// first and the check follows it.
func (c *Cart) Total() Money {
	var total Money

	for _, ci := range c.items {
		total += ci.Item.Price * Money(ci.Quantity)
	}

	return total
}

// Items returns a shallow copy of the cart's lines.
//
// Modifying the returned slice does not change the cart's internal slice.
func (c Cart) Items() []CartItem {
	return slices.Clone(c.items)
}

func (c *Cart) indexOf(id string) int {
	for i, ci := range c.items {
		if ci.Item.ID == id {
			return i
		}
	}
	return -1
}
