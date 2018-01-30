package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Create_Log_With_200_Status_Code(t *testing.T) {
	now := time.Now()
	l := CreateLogger()
	log := Log{
		Time: now,
		Code: 200,
		Hits: 3,
	}
	l.Push(log)
	hits := l.GetByTime(now)[200]
	assert.Equal(t, 3, int(hits))
}
func Test_Create_Log_With_Multiple_Status_Code(t *testing.T) {
	now := time.Now()
	l := CreateLogger()

	var log Log
	log = Log{Time: now, Code: 200, Hits: 3}
	l.Push(log)
	log = Log{Time: now, Code: 200, Hits: 5}
	l.Push(log)
	log = Log{Time: now, Code: 404, Hits: 2}
	l.Push(log)

	var hits int64
	hits = l.GetByTime(now)[200]
	assert.Equal(t, 8, int(hits))

	hits = l.GetByTime(now)[404]
	assert.Equal(t, 2, int(hits))
}
func Test_Create_Log_With_Diffrenet_Times(t *testing.T) {
	time1 := time.Now()
	time2 := time.Now().Add(time.Minute * time.Duration(2))

	l := CreateLogger()
	var log Log
	log = Log{Time: time1, Code: 200, Hits: 3}
	l.Push(log)
	log = Log{Time: time1, Code: 200, Hits: 5}
	l.Push(log)
	log = Log{Time: time2, Code: 200, Hits: 5}
	l.Push(log)
	log = Log{Time: time2, Code: 404, Hits: 2}
	l.Push(log)

	var hits int64
	hits = l.GetByTime(time1)[200]
	assert.Equal(t, 8, int(hits))

	hits = l.GetByTime(time2)[200]
	assert.Equal(t, 5, int(hits))

	hits = l.GetByTime(time1)[404]
	assert.Equal(t, 0, int(hits))

	hits = l.GetByTime(time2)[404]
	assert.Equal(t, 2, int(hits))
}
