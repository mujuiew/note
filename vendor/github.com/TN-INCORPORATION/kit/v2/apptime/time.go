package apptime

import (
	"log"
	"os"
	"time"
)

var timeOffset time.Duration
var systemDate time.Time
var thLocation *time.Location
var isGetFromFile bool
var timeOffsetStr string
var sysdateStr string
var syncTimeStr string
var syncTimeUrl string

func init() {

	var err error
	thLocation, _ = time.LoadLocation("Asia/Bangkok")

	//APP_OFFSET_TIME
	timeOffsetStr = os.Getenv("APP_OFFSET_TIME")
	if timeOffsetStr != "" {
		timeOffset, err = time.ParseDuration(timeOffsetStr)
		if err != nil {
			panic("Override App time Failed with APP_OFFSET_TIME : " + err.Error())
		}
	}

	//APP_SYSTEMDATE
	sysdateStr = os.Getenv("APP_SYSTEMDATE")

	//APP_SYNC_TIME_SERVER
	syncTimeUrl = os.Getenv("APP_SYNC_TIME_SERVER")

	log.Println("APP TIME : ", Now())
}

func Now() time.Time {

	if syncTimeUrl != "" {
		currentDate, err := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), thLocation)
		if err == nil {
			syncTimeStr = GetSyncTime(syncTimeUrl)
			systemDateTime, err := time.ParseInLocation("2006-01-02 15:04", syncTimeStr, thLocation)
			if err == nil {
				return systemDateTime
			}
			systemDate, err := time.ParseInLocation("2006-01-02", syncTimeStr, thLocation)
			if err == nil {
				offsetDate := systemDate.Sub(currentDate)
				return time.Now().Add(offsetDate).In(thLocation)
			}
			log.Println("APP SYNC TIME: call server fail, change to use default time")
			syncTimeUrl = ""
		} else {
			log.Println("APP SYNC TIME: get time fail, change to use default time")
			syncTimeUrl = ""
		}
	}
	if sysdateStr == "GET-FROM-FILE" {
		currentDate, err := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), thLocation)
		if err == nil {
			dateChangeInfo := getDateFromFile()
			var err error
			systemDate, err = time.ParseInLocation("2006-01-02", dateChangeInfo.Date, thLocation)
			if err != nil {
				panic(err.Error())
			}
			offsetDate := systemDate.Sub(currentDate)
			return time.Now().Add(offsetDate).In(thLocation)
		}
		log.Println("APP TIME FILE: get time fail, change to use default time")
		sysdateStr = ""
	}

	if timeOffsetStr != "" {
		return time.Now().Add(timeOffset).In(thLocation)
	}
	return time.Now().In(thLocation)
}

func SetDateTimeForUnitTest(date time.Time) {
	//override time off set with target date - now date = -different duration between target and now date
	timeOffset = date.Sub(time.Now())
	//set offset str for case timeOffsetStr != "" --> time = now time + (-diff duration between target and now date) = target date
	timeOffsetStr = "unit_test"
}
