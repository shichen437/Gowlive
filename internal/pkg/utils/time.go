package utils

import (
	"github.com/gogf/gf/v2/os/gtime"
)

func Now() *gtime.Time {
	return gtime.NewFromStrFormat(gtime.Now().Local().Format("Y-m-d H:i:s"), "Y-m-d H:i:s")
}

func CalcNextDate(cycleNum, cycleType, cycleDay int, startDate *gtime.Time) *gtime.Time {
	if cycleType == 0 {
		return nil
	}

	if startDate == nil {
		startDate = Now()
	}

	now := Now()
	nextDate := startDate.Clone()

	for !nextDate.After(now) {
		currentDate := nextDate.Clone()
		switch cycleType {
		case 1:
			switch cycleDay {
			case 1:
				nextDate = currentDate.AddDate(0, 0, cycleNum*372)
			case 2:
				nextDate = currentDate.AddDate(0, 0, cycleNum*360)
			default:
				nextDate = currentDate.AddDate(cycleNum, 0, 0)
			}
		case 2:
			switch cycleDay {
			case 1:
				nextDate = currentDate.AddDate(0, 0, cycleNum*93)
			case 2:
				nextDate = currentDate.AddDate(0, 0, cycleNum*90)
			default:
				nextDate = currentDate.AddDate(0, cycleNum*3, 0)
			}
		case 3:
			switch cycleDay {
			case 1:
				nextDate = currentDate.AddDate(0, 0, cycleNum*31)
			case 2:
				nextDate = currentDate.AddDate(0, 0, cycleNum*30)
			default:
				nextDate = currentDate.AddDate(0, cycleNum, 0)
			}
		case 4:
			nextDate = currentDate.AddDate(0, 0, cycleNum*7)
		default:
			return nil
		}
	}

	return nextDate
}

func ClacDateBetween(now, nextDate *gtime.Time) int {
	if !now.Before(nextDate) {
		return -99
	}
	d1 := gtime.NewFromStrFormat(now.Format("Y-m-d"), "Y-m-d")
	d2 := gtime.NewFromStrFormat(nextDate.Format("Y-m-d"), "Y-m-d")
	duration := d2.Sub(d1).Hours()
	days := int(duration / (24))

	return days
}
