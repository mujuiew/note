package apptime

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

func getDateFromFile() DateChangeInfo {
	var dateChangeInfo DateChangeInfo
	f, err := os.Open(string(os.PathSeparator) + "tmp" + string(os.PathSeparator) + "appdate.json")
	if err != nil {
		dateChangeInfo.Date = time.Now().Format("2006-01-02")
		return dateChangeInfo
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &dateChangeInfo)
	if err != nil {
		panic(err)
	}
	return dateChangeInfo
}

type DateChangeInfo struct {
	Date string `json:"app_date"`
}
