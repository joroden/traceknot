package dashboard

import "time"

func weekStart(now time.Time) time.Time {
	local := now.Local()
	daysSinceMonday := (int(local.Weekday()) + 6) % 7
	day := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, local.Location())
	return day.AddDate(0, 0, -daysSinceMonday)
}
