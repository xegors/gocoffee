package domain_test

import (
	"testing"

	"github.com/xegors/gocoffee/internal/domain"
)

func TestRankFor(t *testing.T) {
	tests := []struct {
		name  string
		spent domain.Money
		want  domain.Rank
	}{
		{name: "negative spent", spent: -100, want: domain.RankNone},
		{name: "zero spent", spent: 0, want: domain.RankNone},
		{name: "below bronze", spent: 3999, want: domain.RankNone},
		{name: "exactly bronze", spent: 4000, want: domain.RankBronze},
		{name: "below silver", spent: 8999, want: domain.RankBronze},
		{name: "exactly silver", spent: 9000, want: domain.RankSilver},
		{name: "below gold", spent: 19999, want: domain.RankSilver},
		{name: "exactly gold", spent: 20000, want: domain.RankGold},
		{name: "below platinum", spent: 49999, want: domain.RankGold},
		{name: "exactly platinum", spent: 50000, want: domain.RankPlatinum},
		{name: "above platinum", spent: 100000, want: domain.RankPlatinum},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.RankFor(tt.spent)

			if got != tt.want {
				t.Errorf("| got %2d | want %2d |", got, tt.want)
			}
		})
	}
}
