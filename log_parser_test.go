package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Parse_Log(t *testing.T) {
	logLine := "[2018-01-30T18:42:51+02:00] \"GET / HTTP/2.0\" 200"
	parser := parseLog(logLine)
	assert.Equal(t, 200, parser.Code)
	assert.Equal(t, "/", parser.Path)
	expectedTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2018-01-30T18:42:51+02:00")
	assert.Equal(t, expectedTime, parser.Time)

	logLine2 := "[2018-01-30T18:32:12+02:00] \"GET /asdasd HTTP/2.0\" 404"
	parser2 := parseLog(logLine2)
	assert.Equal(t, 404, parser2.Code)
	assert.Equal(t, "/asdasd", parser2.Path)
	expectedTime2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2018-01-30T18:32:12+02:00")
	assert.Equal(t, expectedTime2, parser2.Time)
}
