package eel

import (
	"github.com/gingerxman/eel/log"
	"time"
	"strings"
)

var locShanghai, _ = time.LoadLocation("Asia/Shanghai")

func ParseTime(strTime string) time.Time {
	if strings.Count(strTime, ":") == 1 {
		strTime += ":00"
	}
	
	timeVal, err := time.ParseInLocation("2006-01-02 15:04:05", strTime, locShanghai)
	if err != nil {
		log.Logger.Error(err)
	}
	return timeVal
}

