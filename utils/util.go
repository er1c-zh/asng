package utils

import "time"

func GetTodayWithOffset(h, M, s int) time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, h, M, s, 0, time.Local)
}
