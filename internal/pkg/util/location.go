package util

import "time"

var Location *time.Location

func init() {
	Location = time.Now().Location()
	asiaTaipei, err := time.LoadLocation("Asia/Taipei")
	if err == nil {
		Location = asiaTaipei
	}
}
