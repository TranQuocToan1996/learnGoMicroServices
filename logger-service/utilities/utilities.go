package utilities

import "time"

var (
	HoChiMinhCityTimeZone, _ = time.LoadLocation("Asia/Saigon")

	TimeLocalNow func() time.Time = func() time.Time {
		return timeLocalNow(HoChiMinhCityTimeZone)
	}
)

func timeLocalNow(local *time.Location) time.Time {
	return time.Now().In(local)
}
