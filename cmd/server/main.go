package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/czertbytes/pixelizer/handler"
)

func main() {
	httpPort := flag.Int("httport", 8080, "HTTP port where web service runs.")
	flag.Parse()

	go func() {
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	httpCh := make(chan struct{})
	go func() {
		defer close(httpCh)

		addr := fmt.Sprintf(":%d", *httpPort)
		http.Handle("/", &handler.Index{})
		http.Handle("/pixelize", &handler.Pixelize{})
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal("failed to create HTTP listener")
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-sigCh:
		log.Println("Signal received, terminating")
	case <-httpCh:
		log.Println("HTTP listener finished, terminating")
	}
}
