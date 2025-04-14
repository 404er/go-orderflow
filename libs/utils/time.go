package utils

import "time"

func GetDiffDays(start, end time.Time) int {
	startTime := time.Unix(start.Unix(), 0)
	endTime := time.Unix(end.Unix(), 0)

	diff := endTime.Sub(startTime)
	return int(diff.Hours() / 24)
}
