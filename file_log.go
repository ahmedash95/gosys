package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"
)

var mutex sync.Mutex
var entries logEntries

var logPath = "/var/log/gosys.log"

var tickCh = time.Tick(2 * time.Second)
var writeDelay = 2 * time.Second

type LogEntry struct {
	Timestamp time.Time
	Message   string
}

type logEntries []LogEntry

func (le logEntries) Len() int {
	return len(le)
}

func (le logEntries) Swap(i, j int) {
	le[i], le[j] = le[j], le[i]
}

func (le logEntries) Less(i, j int) bool {
	return le[i].Timestamp.Before(le[j].Timestamp)
}

func LogToFile(e LogEntry) {
	mutex.Lock()
	entries = append(entries, e)
	mutex.Unlock()
}

func SetLogPath(path string) {
	logPath = path
}

func writeLog() {
	for range tickCh {
		mutex.Lock()

		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			fmt.Println(err)
			// detect if file exists
			var _, err = os.Stat(logPath)
			// create file if not exists
			if os.IsNotExist(err) {
				var file, err = os.Create(logPath)
				if err != nil {
					fmt.Println(err)
				}
				defer file.Close()
			}
			mutex.Unlock()
			continue
		}
		targetTime := time.Now().Add(-writeDelay)
		sort.Sort(entries)
		for i, entry := range entries {
			if entry.Timestamp.Before(targetTime) {
				_, err := logFile.WriteString(writeEntry(entry))
				if err != nil {
					fmt.Println(err)
				}

				if i == len(entries)-1 {
					entries = logEntries{}
				}

			} else {
				entries = entries[i:]
				break
			}
		}

		logFile.Close()

		mutex.Unlock()
	}
}

func writeEntry(entry LogEntry) string {
	return fmt.Sprintf("%v\n", entry.Message)
}

func readLogs() ([]string, error) {
	var lines []string
	file, err := ioutil.ReadFile(logPath)
	if err != nil {
		return lines, err
	}
	buf := bytes.NewBuffer(file)
	for {
		line, err := buf.ReadString('\n')
		if len(line) == 0 {
			if err != nil {
				if err == io.EOF {
					break
				}
				return lines, err
			}
		}
		lines = append(lines, line)
		if err != nil && err != io.EOF {
			return lines, err
		}
	}
	return lines, nil
}
