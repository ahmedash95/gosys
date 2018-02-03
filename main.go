package main

import (
	"flag"
	"fmt"
	"net/http"

	syslog "gopkg.in/mcuadros/go-syslog.v2"
)

var (
	Logs     Logger
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

	// init logger
	Logs = CreateLogger()

	go func() {
		server.Wait()
	}()

	fmt.Printf("Server is up and running on %s\npress <Enter> key to exit\n", *HTTPAddr)

	// start websocket listner
	go handleMessages()

	go func() {

		// Configure websocket route
		http.HandleFunc("/ws", handleConnections)
		// web ui dashboard
		fs := http.FileServer(http.Dir("ui"))
		http.Handle("/", fs)
		http.ListenAndServe(":3000", nil)
	}()

	fmt.Scanln()
}

func pushLog(line string) {
	l := parseLog(line)
	log := Log{
		Time: l.Time,
		Code: l.Code,
		Hits: 1,
	}
	Logs.Push(log)
	// Send the newly received message to the broadcast channel
	broadcast <- log
	// clearCli()
	// fmt.Println("Push new log")
	// fmt.Println(Logs)
}
