package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var topic chan string

func sse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	flusher, _ := w.(http.Flusher)

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	// go func() {
	for {
		select {
		case v := <-topic:
			fmt.Fprintf(w, "receive: %s\n", v)
			log.Println("hey")
			flusher.Flush()
		case <-r.Context().Done():
			log.Println("connection closed")
			return
		}
	}
	// }()
}

func main() {
	topic = make(chan string)
	go func() {
		for range time.Tick(time.Second) {
			topic <- time.Now().Format("2006-01-02 15:04:05")
		}
	}()

	http.HandleFunc("/sse", sse)
	http.ListenAndServe(":8080", nil)
}
