package domain

// Rank is a User's personal discount level, stored as a percentage.
//
// A rank never decreases: total spent money only grows, so RankFor is a pure function
// of that total. RankNone means no discount yet.
type Rank int

// Rank values. Each value is the discount percentage it grants.
const (
	RankNone     Rank = 0
	RankBronze   Rank = 3
	RankSilver   Rank = 5
	RankGold     Rank = 7
	RankPlatinum Rank = 10
)

// Threshold for reaching each rank, in rubles in total spent money.
const (
	bronzeThreshold   Money = 4000
	silverThreshold   Money = 9000
	goldThreshold     Money = 20000
	platinumThreshold Money = 50000
)

// RankFor returns a Rank granted by total spent money.
//
// Thresholds are checked from highest to lowest because the first match wins.
// Spending below bronzeThreshold yields RankNone, which also covers negative totals.
func RankFor(ts Money) Rank {
	switch {
	case ts >= platinumThreshold:
		return RankPlatinum
	case ts >= goldThreshold:
		return RankGold
	case ts >= silverThreshold:
		return RankSilver
	case ts >= bronzeThreshold:
		return RankBronze

	default:
		return RankNone
	}
}
