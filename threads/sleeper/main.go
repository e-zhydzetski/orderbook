package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func errorAwareMain() error {
	var listenAddr string

	flag.StringVar(&listenAddr, "listen-addr", ":8080", "listening address")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		delay := time.Second
		if delayParam := r.URL.Query().Get("delay"); delayParam != "" {
			if d, err := time.ParseDuration(delayParam); err == nil {
				delay = d
			}
		}
		time.Sleep(delay)
		_, _ = w.Write([]byte(delay.String()))
	})

	return http.ListenAndServe(listenAddr, nil)
}

func main() {
	if err := errorAwareMain(); err != nil {
		log.Print(err)
	}
}
