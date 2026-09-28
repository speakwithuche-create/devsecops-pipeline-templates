package retention

import (
	"testing"
	"time"
)

func daysBack(n int) []time.Time {
	start := time.Date(2026, 9, 28, 2, 0, 0, 0, time.UTC)
	var out []time.Time
	for i := 0; i < n; i++ {
		out = append(out, start.AddDate(0, 0, -i))
	}
	return out
}

func TestKeepsLastSevenDailies(t *testing.T) {
	got := Policy{Daily: 7}.Keep(daysBack(30))
	if len(got) != 7 {
		t.Fatalf("want 7, got %d", len(got))
	}
	if !got[0].Equal(daysBack(1)[0]) {
		t.Fatal("newest point must be kept first")
	}
}

func TestWeeklyAndMonthlyExtendHistory(t *testing.T) {
	got := Policy{Daily: 7, Weekly: 4, Monthly: 3}.Keep(daysBack(120))
	// 3 monthly points: the newest of September, August and July.
	if oldest := got[len(got)-1]; oldest.Month() != time.July || oldest.Day() != 31 {
		t.Fatalf("oldest kept point should be July 31, got %v", oldest)
	}
	// 7 dailies + older weeklies + older monthlies, with no duplicates.
	seen := map[time.Time]bool{}
	for _, p := range got {
		if seen[p] {
			t.Fatalf("duplicate point %v", p)
		}
		seen[p] = true
	}
}

func TestEmptyPolicyKeepsNothing(t *testing.T) {
	if got := (Policy{}).Keep(daysBack(10)); len(got) != 0 {
		t.Fatalf("want 0, got %d", len(got))
	}
}

func TestDoesNotMutateInput(t *testing.T) {
	in := daysBack(5)
	in[0], in[4] = in[4], in[0]
	first := in[0]
	Policy{Daily: 2}.Keep(in)
	if !in[0].Equal(first) {
		t.Fatal("Keep must not reorder the caller's slice")
	}
}
