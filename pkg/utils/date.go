package utils

import "time"

func CalcExpiredAt(days uint32) time.Time {
	return time.Now().Add(time.Duration(days) * time.Hour * 24)
}
