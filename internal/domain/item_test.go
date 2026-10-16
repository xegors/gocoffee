package domain_test

import (
	"errors"
	"testing"

	"github.com/xegors/gocoffee/internal/domain"
)

func TestParseCategory(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    domain.Category
		wantErr error
	}{
		{name: "coffee with trailing spaces", input: "  coffee\t ", want: domain.CategoryCoffee},
		{name: "capitalized accessories", input: "Accessories", want: domain.CategoryAccessories},
		{name: "sets", input: "sets", want: domain.CategorySets},
		{name: "surprise", input: "surprise", want: domain.CategorySurprise},
		{name: "merch deferred on stage 6", input: "merch", wantErr: domain.ErrUnknownCategory},
		{name: "empty category", input: "", wantErr: domain.ErrUnknownCategory},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ParseCategory(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
