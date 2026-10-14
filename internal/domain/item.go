// Package domain holds the entities and rules of the business domain.
//
// The package knows nothing about persistence: encoding/json never reaches
// it. The standard library is allowed, but only for computation, never for
// the storage format.
package domain

import (
	"strings"
)

// Category is the closed set of catalog categories.
//
// Values are normalized by ParseCategory: surrounding whitespace is
// trimmed and the case is lowered. A stored value is therefore always
// lowercase, so comparing two categories with == is safe. A literal like
// Category("Coffee") bypasses that rule.
type Category string

// Category values. Latin script, because these strings land in items.json.
const (
	CategoryCoffee      Category = "coffee"
	CategoryAccessories Category = "accessories"
	CategorySets        Category = "sets"
	CategorySurprise    Category = "surprise"
)

// ParseCategory normalizes s and checks it against the set: it trims
// whitespace, lowers the case, then matches.
//
// It returns ErrUnknownCategory if the value is not recognized.
func ParseCategory(s string) (Category, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	switch Category(s) {
	case CategoryCoffee:
		return CategoryCoffee, nil
	case CategoryAccessories:
		return CategoryAccessories, nil
	case CategorySets:
		return CategorySets, nil
	case CategorySurprise:
		return CategorySurprise, nil

	default:
		return "", ErrUnknownCategory
	}
}

// BrewMethod is how a bean is roasted. It is not a category: an Ethiopian
// Yirgacheffe goes to filter, a Brazilian Santos goes to espresso, which
// makes it a property of the item rather than a kind of item.
//
// Unlike Category, the empty value is valid and meaningful: BrewNone means
// the field does not apply, because the item is not a bean.
type BrewMethod string

// BrewMethod values. BrewNone is the empty string on purpose, so that an
// item with no roast method stays distinguishable from an item carrying a
// typo.
const (
	BrewNone     BrewMethod = ""
	BrewFilter   BrewMethod = "filter"
	BrewEspresso BrewMethod = "espresso"
)

// ParseBrewMethod behaves like ParseCategory: it trims whitespace, lowers
// the case and matches against the set. An empty string resolves to
// BrewNone.
//
// It returns ErrUnknownBrew if the value is not recognized.
func ParseBrewMethod(s string) (BrewMethod, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	switch BrewMethod(s) {
	case BrewNone:
		return BrewNone, nil
	case BrewFilter:
		return BrewFilter, nil
	case BrewEspresso:
		return BrewEspresso, nil

	default:
		return "", ErrUnknownBrew
	}
}

// Item is a catalog entry: the thing that is sold.
//
// The struct is flat even though some fields make sense only for part of the
// categories: WeightGrams and Brew are filled in for beans alone and stay at
// zero for an accessory or a set. That debt is deliberate. Splitting per
// category turns the catalog polymorphic and forces a type discriminator
// into JSON, so the split is deferred to Stage 7.
type Item struct {
	ID          string
	Name        string
	Category    Category
	Price       Money
	WeightGrams int        // 0 means not applicable: a bean has weight, a dripper does not
	Brew        BrewMethod // BrewNone means not applicable
}
