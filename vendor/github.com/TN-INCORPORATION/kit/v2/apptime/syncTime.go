package apptime

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type SyncDate struct {
	Date string `json:"date" `
}

func GetSyncTime(url string) string {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	var syncDate SyncDate
	err = json.Unmarshal(body, &syncDate)
	if err != nil {
		return ""
	}
	return syncDate.Date
}
