package timezone

import "time"

func init() {

	thLocation, _ = time.LoadLocation("Asia/Bangkok")
}

var thLocation *time.Location

func GetTimeZone() *time.Location {
	return thLocation
}
