package domain

import (
	"math"
	"testing"
)

// Тест внутри пакета, потому что addMoney неэкспортируема: снаружи её не видно,
// и проверять её через Deposit — это проверять не то правило, а обёртку над ним.
func TestAddMoney(t *testing.T) {
	tests := []struct {
		name   string
		a      Money
		b      Money
		want   Money
		wantOK bool
	}{
		{name: "zeros", a: 0, b: 0, want: 0, wantOK: true},
		{name: "positive amounts", a: 1000, b: 200, want: 1200, wantOK: true},
		{name: "zero left operand", a: 0, b: 200, want: 200, wantOK: true},
		{name: "negative left operand", a: -100, b: 200, wantOK: false},
		{name: "negative right operand", a: 100, b: -200, wantOK: false},
		{name: "sum at the money limit", a: math.MaxInt64 - 1, b: 1, want: math.MaxInt64, wantOK: true},
		{name: "sum over the money limit", a: math.MaxInt64, b: 1, wantOK: false},
		{name: "limit plus zero", a: math.MaxInt64, b: 0, want: math.MaxInt64, wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Флаг проверяется отдельно от значения: addMoney на отказе
			// возвращает 0, и тест, глядящий только на значение, прошёл бы
			// на любом отказе с любым числом.
			got, ok := addMoney(tt.a, tt.b)
			if ok != tt.wantOK {
				t.Fatalf("got ok %t, want %t", ok, tt.wantOK)
			}

			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
