package main

import (
	"regexp"
	"strconv"
	"time"
)

type LogPars struct {
	Time time.Time
	Path string
	Code int
}

func parseLog(line string) LogPars {
	re := regexp.MustCompile("\\[(.*?)\\]\\s\\\"(.*?)\\s(.*?)\\s(.*)\\\"\\s(.*?)$")
	q := re.FindStringSubmatch(line)
	date, _ := time.Parse("2006-01-02T15:04:05Z07:00", q[1])
	StatusCode, _ := strconv.Atoi(q[5])
	log := LogPars{
		Time: date,
		Path: q[3],
		Code: StatusCode,
	}
	return log
}
