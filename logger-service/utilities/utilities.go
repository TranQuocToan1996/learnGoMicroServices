package utilities

import "time"

var (
	HoChiMinhCityTimeZone, _ = time.LoadLocation("Asia/Saigon")

	// TODO: add saigon location
	TimeLocalNow func() time.Time = func() time.Time {
		return timeLocalNow(time.Now().Location())
	}
)

func timeLocalNow(local *time.Location) time.Time {
	return time.Now().In(local)
}
