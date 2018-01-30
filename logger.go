package main

import "time"

var defaultTimeFormat = "2006-01-02 15:04"

type Logger struct {
	Logs       map[string]map[int]int64
	TimeFormat string
}

type Log struct {
	Time time.Time
	Code int
	Hits int64
}

func setDefaultTimeFormat(t string) {
	defaultTimeFormat = t
}

func getDefaultTimeFormat() string {
	return defaultTimeFormat
}

func (l *Logger) setTimeFormat(f string) {
	l.TimeFormat = f
}

func CreateLogger() Logger {
	l := Logger{
		Logs:       make(map[string]map[int]int64),
		TimeFormat: getDefaultTimeFormat(),
	}
	return l
}

func (l *Logger) Push(log Log) {
	k := log.Time.Format(l.TimeFormat)
	logTime, ok := l.Logs[k]
	if !ok {
		l.Logs[k] = make(map[int]int64)
		logTime = l.Logs[k]
		logTime[log.Code] = 0
	}
	logTime[log.Code] += log.Hits
}

func (l *Logger) GetByTime(t time.Time) map[int]int64 {
	return l.Get()[t.Format(l.TimeFormat)]
}

func (l *Logger) Get() map[string]map[int]int64 {
	return l.Logs
}
