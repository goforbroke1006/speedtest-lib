package main

import (
	"log"
	"net/http"
	"time"

	"github.com/goforbroke1006/speedtest-lib/handler/test_speed"
	"github.com/goforbroke1006/speedtest-lib/loader"
	"github.com/goforbroke1006/speedtest-lib/upgrader"
)

func main() {
	ooklaUpgrader := upgrader.NewUpgrader(loader.NewOoklaLoader(), time.Minute)
	go ooklaUpgrader.Run()

	netflixUpgrader := upgrader.NewUpgrader(loader.NewNetflixLoader(), time.Minute)
	go netflixUpgrader.Run()

	sources := map[string]upgrader.Upgrader{
		"ookla":   ooklaUpgrader,
		"netflix": netflixUpgrader,
	}

	http.HandleFunc("/healthz", healthHandlerMiddleware())
	http.HandleFunc("/readyz", readyHandlerMiddleware(sources))
	http.HandleFunc("/test", test_speed.TestSpeedHandler(sources))

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func healthHandlerMiddleware() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

func readyHandlerMiddleware(sources map[string]upgrader.Upgrader) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, s := range sources {
			if !s.IsReady() {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte("fail"))
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
