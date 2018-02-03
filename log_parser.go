package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogPars struct {
	Time time.Time
	Path string
	Code int
}

func parseLog(line string) LogPars {
	line = strings.TrimSpace(line)
	if line == "" {
		return LogPars{}
	}
	re := regexp.MustCompile("\\[(.*?)\\]\\s\\\"(.*?)\\s(.*?)\\s(.*)\\\"\\s(.*?)$")
	q := re.FindStringSubmatch(line)
	if len(q) < 1 {
		panic(fmt.Sprintf("Faild to parse log message : %s\n", line))
	}
	date, _ := time.Parse("2006-01-02T15:04:05Z07:00", q[1])
	StatusCode, _ := strconv.Atoi(q[5])
	log := LogPars{
		Time: date,
		Path: q[3],
		Code: StatusCode,
	}
	return log
}
