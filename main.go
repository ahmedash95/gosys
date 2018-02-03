package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	syslog "gopkg.in/mcuadros/go-syslog.v2"
)

var (
	Logs     Logger
	OldLogs  []websocketLogMessage
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "0.0.0.0:12301", "Address to listen for UDP requests on")
)

func main() {

	flag.Parse()

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)
	server.ListenUDP(*HTTPAddr)
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			work := WorkRequest{Line: logParts["content"].(string)}
			// Push the work onto the queue.
			WorkQueue <- work
		}
	}(channel)

	StartDispatcher(*NWorkers)

	// Register log path
	SetLogPath(getProjectPath() + "/gosys.log")

	// init logger
	Logs = CreateLogger()

	go func() {
		server.Wait()
	}()

	fmt.Printf("Server is up and running on %s\npress <Enter> key to exit\n", *HTTPAddr)

	// start websocket listner
	go handleMessages()

	go writeLog()

	go func() {
		// Configure websocket route
		http.HandleFunc("/ws", handleConnections)
		// web ui dashboard
		fs := http.FileServer(http.Dir("ui"))
		http.Handle("/", fs)
		http.HandleFunc("/logs", loadOldLogs)
		http.ListenAndServe(":3000", nil)
	}()

	go func() {
		lines, err := readLogs()
		if err != nil {
			fmt.Println("Can't load old logs from file")
			return
		}
		for _, l := range lines {
			p := parseLog(l)
			wsm := websocketLogMessage{
				Time: p.Time.Format(websocketTimeFormat),
				Hits: 1,
				Code: p.Code,
			}
			OldLogs = append(OldLogs, wsm)
		}
	}()

	fmt.Scanln()
}

func loadOldLogs(w http.ResponseWriter, r *http.Request) {
	var msgs = OldLogs
	if len(OldLogs) > 3000 {
		msgs = OldLogs[len(OldLogs)-3000:]
	}
	body, _ := json.Marshal(msgs)
	w.Write(body)
}

func pushLog(line string) {
	l := parseLog(line)
	log := Log{
		Time: l.Time,
		Code: l.Code,
		Hits: 1,
	}
	// Store Log in Memory
	Logs.Push(log)
	// Send the newly received message to the broadcast channel
	broadcast <- log
	// Log to file
	e := LogEntry{
		Timestamp: log.Time,
		Message:   line,
	}
	LogToFile(e)
}
