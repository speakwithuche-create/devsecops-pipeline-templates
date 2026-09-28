// Package retention decides which backup restore points to keep under a
// grandfather-father-son policy.
package retention

import (
	"fmt"
	"sort"
	"time"
)

// Policy keeps the newest point of each of the last Daily days, Weekly ISO
// weeks and Monthly months.
type Policy struct {
	Daily, Weekly, Monthly int
}

// Keep returns the restore points to retain, newest first.
func (p Policy) Keep(points []time.Time) []time.Time {
	sorted := append([]time.Time(nil), points...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].After(sorted[j]) })

	keep := map[time.Time]bool{}
	pick := func(limit int, bucket func(time.Time) string) {
		seen := map[string]bool{}
		for _, t := range sorted {
			if len(seen) == limit {
				return
			}
			if b := bucket(t); !seen[b] {
				seen[b] = true
				keep[t] = true
			}
		}
	}
	pick(p.Daily, func(t time.Time) string { return t.UTC().Format("2006-01-02") })
	pick(p.Weekly, func(t time.Time) string {
		y, w := t.UTC().ISOWeek()
		return fmt.Sprintf("%d-W%02d", y, w)
	})
	pick(p.Monthly, func(t time.Time) string { return t.UTC().Format("2006-01") })

	var out []time.Time
	for _, t := range sorted {
		if keep[t] {
			out = append(out, t)
		}
	}
	return out
}
