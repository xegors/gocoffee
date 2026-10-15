package domain_test

import (
	"errors"
	"math"
	"testing"

	"github.com/xegors/gocoffee/internal/domain"
)

func item(id string, price domain.Money) domain.Item {
	return domain.Item{
		ID:    id,
		Price: price,
	}
}

func mustAdd(t *testing.T, c *domain.Cart, id string, price domain.Money, qty int) {
	t.Helper()
	if err := c.Add(domain.Item{ID: id, Price: price}, qty); err != nil {
		t.Fatalf("setup: cannot add %q: %v", id, err)
	}
}

func TestCartAdd(t *testing.T) {
	tests := []struct {
		name      string
		initial   []domain.CartItem // что уже в корзине
		input     domain.Item
		qty       int
		wantLines int // сколько строк ожидаем
		wantQty   int // количество в найденной строке
		wantErr   error
	}{
		{name: "into empty cart", input: item("ethiopia", 799), qty: 2, wantLines: 1, wantQty: 2},
		{name: "increase existing",
			initial: []domain.CartItem{{Item: item("ethiopia", 799), Quantity: 1}},
			input:   item("ethiopia", 799), qty: 3,
			wantLines: 1, wantQty: 4,
		},
		{name: "new item append",
			initial: []domain.CartItem{{Item: item("ethiopia", 799), Quantity: 1}},
			input:   item("kenya", 999), qty: 2,
			wantLines: 2, wantQty: 2,
		},
		{name: "reject empty id", input: domain.Item{}, qty: 1,
			wantLines: 0, wantQty: 0, wantErr: domain.ErrInvalidItem,
		},
		{name: "reject zero qty", input: item("ethiopia", 799), qty: 0,
			wantLines: 0, wantQty: 0, wantErr: domain.ErrInvalidQuantity,
		},
		{name: "reject negative qty", input: item("ethiopia", 799), qty: -2,
			wantLines: 0, wantQty: 0, wantErr: domain.ErrInvalidQuantity,
		},
		{name: "reject overflow",
			initial: []domain.CartItem{{Item: item("kenya", 999), Quantity: math.MaxInt}},
			input:   item("kenya", 999), qty: 1,
			wantLines: 1, wantQty: math.MaxInt, wantErr: domain.ErrInvalidQuantity,
		},
		{name: "accept exact max",
			initial: []domain.CartItem{{Item: item("kenya", 999), Quantity: math.MaxInt - 1}},
			input:   item("kenya", 999), qty: 1,
			wantLines: 1, wantQty: math.MaxInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &domain.Cart{}
			for _, line := range tt.initial {
				mustAdd(t, c, line.Item.ID, line.Item.Price, line.Quantity)
			}

			err := c.Add(tt.input, tt.qty)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %q, want %q", err, tt.wantErr)
			}

			lines := c.Items()
			if len(lines) != tt.wantLines {
				t.Errorf("got %d lines, want %d lines", len(lines), tt.wantLines)
			}
			found := false
			for _, line := range lines {
				if line.Item.ID == tt.input.ID {
					found = true
					if line.Quantity != tt.wantQty {
						t.Errorf("got qty %d, want qty %d", line.Quantity, tt.wantQty)
					}
				}
			}
			if !found && tt.wantErr == nil {
				t.Errorf("item %q not found in cart", tt.input.ID)
			}
		})
	}
}

func TestCartRemove(t *testing.T) {
	tests := []struct {
		name      string
		initial   []domain.CartItem // что уже в корзине
		removeID  string
		removeQty int
		wantLines int      // сколько строк ожидаем
		wantIDs   []string // порядок идентификаторов после операции
		wantID    string   // какую строку проверяем количеством, "" = никакую
		wantQty   int      // количество в строке wantID
		wantErr   error
	}{
		{name: "partial",
			initial:  []domain.CartItem{{Item: item("a", 1), Quantity: 5}},
			removeID: "a", removeQty: 2,
			wantLines: 1, wantIDs: []string{"a"}, wantID: "a", wantQty: 3,
		},
		{name: "exact remove line",
			initial:  []domain.CartItem{{Item: item("a", 1), Quantity: 5}},
			removeID: "a", removeQty: 5,
			wantLines: 0, wantIDs: []string{},
		},
		{name: "more than available",
			initial:  []domain.CartItem{{Item: item("a", 1), Quantity: 5}},
			removeID: "a", removeQty: 7,
			wantLines: 1, wantIDs: []string{"a"}, wantID: "a", wantQty: 5,
			wantErr: domain.ErrInsufficientQuantity,
		},
		{name: "reject zero qty",
			initial:  []domain.CartItem{{Item: item("a", 1), Quantity: 5}},
			removeID: "a", removeQty: 0,
			wantLines: 1, wantIDs: []string{"a"}, wantID: "a", wantQty: 5,
			wantErr: domain.ErrInvalidQuantity,
		},
		{name: "reject negative qty",
			initial:  []domain.CartItem{{Item: item("a", 1), Quantity: 5}},
			removeID: "a", removeQty: -2,
			wantLines: 1, wantIDs: []string{"a"}, wantID: "a", wantQty: 5,
			wantErr: domain.ErrInvalidQuantity,
		},
		{name: "unknown ID",
			removeID: "a", removeQty: 1,
			wantLines: 0, wantIDs: []string{},
			wantErr: domain.ErrItemNotInCart,
		},
		{name: "other line untouched",
			initial: []domain.CartItem{
				{Item: item("a", 1), Quantity: 5},
				{Item: item("b", 3), Quantity: 3},
			},
			removeID: "a", removeQty: 5,
			wantLines: 1, wantIDs: []string{"b"}, wantID: "b", wantQty: 3,
		},
		{name: "middle line keeps order",
			initial: []domain.CartItem{
				{Item: item("a", 1), Quantity: 1},
				{Item: item("b", 3), Quantity: 1},
				{Item: item("c", 5), Quantity: 1},
			},
			removeID: "b", removeQty: 1,
			wantLines: 2, wantIDs: []string{"a", "c"}, wantID: "c", wantQty: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &domain.Cart{}
			for _, line := range tt.initial {
				mustAdd(t, c, line.Item.ID, line.Item.Price, line.Quantity)
			}

			err := c.Remove(tt.removeID, tt.removeQty)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %q, want %q", err, tt.wantErr)
			}

			lines := c.Items()
			if len(lines) != tt.wantLines {
				t.Errorf("got %d lines, want %d lines", len(lines), tt.wantLines)
			}

			gotIDs := make([]string, 0, len(lines))
			for _, line := range lines {
				gotIDs = append(gotIDs, line.Item.ID)
			}
			if len(gotIDs) != len(tt.wantIDs) {
				t.Fatalf("got ids %v, want %v", gotIDs, tt.wantIDs)
			}
			for i, want := range tt.wantIDs {
				if gotIDs[i] != want {
					t.Errorf("position %d: got id %q, want %q", i, gotIDs[i], want)
				}
			}

			if tt.wantID == "" {
				return
			}
			for _, line := range lines {
				if line.Item.ID == tt.wantID {
					if line.Quantity != tt.wantQty {
						t.Errorf("got qty %d for %q, want %d",
							line.Quantity, tt.wantID, tt.wantQty)
					}
					return
				}
			}
			t.Errorf("item %q not found in cart", tt.wantID)
		})
	}
}

func TestCartRemoveAlreadyRemoved(t *testing.T) {
	c := &domain.Cart{}
	mustAdd(t, c, "a", 100, 5)

	if err := c.Remove("a", 5); err != nil {
		t.Fatalf("first remove: %v", err)
	}
	if got := len(c.Items()); got != 0 {
		t.Errorf("after full removal got %d lines, want 0", got)
	}

	err := c.Remove("a", 5)
	if !errors.Is(err, domain.ErrItemNotInCart) {
		t.Errorf("second remove: got error %v, want ErrItemNotInCart", err)
	}
}

func TestCartTotal(t *testing.T) {
	tests := []struct {
		name    string
		initial []domain.CartItem
		want    domain.Money
	}{
		{name: "empty cart", initial: nil, want: 0},
		{name: "single line", initial: []domain.CartItem{
			{Item: item("a", 100), Quantity: 3}}, want: 300},
		{name: "multiple lines", initial: []domain.CartItem{
			{Item: item("a", 100), Quantity: 2},
			{Item: item("b", 250), Quantity: 1}}, want: 450},
		{name: "after dedup", initial: []domain.CartItem{
			{Item: item("a", 100), Quantity: 5}}, want: 500},
		{name: "after partial remove", initial: []domain.CartItem{
			{Item: item("a", 100), Quantity: 3}}, want: 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &domain.Cart{}
			for _, line := range tt.initial {
				mustAdd(t, c, line.Item.ID, line.Item.Price, line.Quantity)
			}

			if got := c.Total(); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCartMerge(t *testing.T) {
	tests := []struct {
		name      string
		left      []domain.CartItem
		right     []domain.CartItem
		wantIDs   []string       // порядок строк после мержа
		wantQty   map[string]int // ID -> количество
		wantLines int
	}{
		{name: "both empty",
			left:      nil,
			right:     nil,
			wantLines: 0,
		},
		{name: "empty left",
			left:      nil,
			right:     []domain.CartItem{{Item: item("a", 5), Quantity: 2}},
			wantIDs:   []string{"a"},
			wantQty:   map[string]int{"a": 2},
			wantLines: 1,
		},
		{name: "empty right",
			left:      []domain.CartItem{{Item: item("a", 5), Quantity: 2}},
			right:     nil,
			wantIDs:   []string{"a"},
			wantQty:   map[string]int{"a": 2},
			wantLines: 1,
		},
		{name: "disjoint",
			left:      []domain.CartItem{{Item: item("a", 5), Quantity: 1}},
			right:     []domain.CartItem{{Item: item("b", 3), Quantity: 1}},
			wantIDs:   []string{"a", "b"},
			wantQty:   map[string]int{"a": 1, "b": 1},
			wantLines: 2,
		},
		{name: "overlaping",
			left:      []domain.CartItem{{Item: item("a", 5), Quantity: 2}},
			right:     []domain.CartItem{{Item: item("a", 3), Quantity: 3}},
			wantIDs:   []string{"a"},
			wantQty:   map[string]int{"a": 5},
			wantLines: 1,
		},
		{name: "order preserved",
			left: []domain.CartItem{
				{Item: item("a", 5), Quantity: 1},
				{Item: item("b", 3), Quantity: 1},
			},
			right:     []domain.CartItem{{Item: item("c", 1), Quantity: 1}},
			wantIDs:   []string{"a", "b", "c"},
			wantQty:   map[string]int{"a": 1, "b": 1, "c": 1},
			wantLines: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := &domain.Cart{}
			for _, line := range tt.left {
				mustAdd(t, left, line.Item.ID, line.Item.Price, line.Quantity)
			}
			right := &domain.Cart{}
			for _, line := range tt.right {
				mustAdd(t, right, line.Item.ID, line.Item.Price, line.Quantity)
			}

			if err := left.Merge(*right); err != nil {
				t.Fatalf("merge: %v", err)
			}

			lines := left.Items()
			if len(lines) != tt.wantLines {
				t.Errorf("got %d lines, want %d", len(lines), tt.wantLines)
			}

			// порядок
			gotIDs := make([]string, 0, len(lines))
			for _, line := range lines {
				gotIDs = append(gotIDs, line.Item.ID)
			}
			if len(gotIDs) != len(tt.wantIDs) {
				t.Fatalf("got ids %v, want %v", gotIDs, tt.wantIDs)
			}
			for i, want := range tt.wantIDs {
				if gotIDs[i] != want {
					t.Errorf("position %d: got %q, want %q", i, gotIDs[i], want)
				}
			}

			for id, wantQty := range tt.wantQty {
				found := false
				for _, line := range lines {
					if line.Item.ID == id {
						found = true
						if line.Quantity != wantQty {
							t.Errorf("got qty %d for %q, want %d", line.Quantity, id, wantQty)
						}
					}
				}
				if !found {
					t.Errorf("item %q not found in cart", id)
				}
			}
		})
	}
}

func TestCartMergePreservesTotal(t *testing.T) {
	pairs := []struct {
		name  string
		left  []domain.CartItem
		right []domain.CartItem
	}{
		{name: "disjoint", left: []domain.CartItem{
			{Item: item("a", 100), Quantity: 2}},
			right: []domain.CartItem{{Item: item("b", 250), Quantity: 1}}},
		{name: "overlapping", left: []domain.CartItem{
			{Item: item("a", 100), Quantity: 2}},
			right: []domain.CartItem{{Item: item("a", 100), Quantity: 3}}},
		{name: "empty right", left: []domain.CartItem{
			{Item: item("a", 100), Quantity: 1}}, right: nil},
	}

	for _, p := range pairs {
		t.Run(p.name, func(t *testing.T) {
			build := func(lines []domain.CartItem) *domain.Cart {
				c := &domain.Cart{}
				for _, line := range lines {
					mustAdd(t, c, line.Item.ID, line.Item.Price, line.Quantity)
				}
				return c
			}

			left, right := build(p.left), build(p.right)
			before := left.Total() + right.Total()

			if err := left.Merge(*right); err != nil {
				t.Fatalf("merge: %v", err)
			}

			if got := left.Total(); got != before {
				t.Errorf("total after merge got %d, want %d", got, before)
			}
		})
	}
}
