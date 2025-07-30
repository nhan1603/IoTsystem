package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"os"
	"syscall"
	"time"
)

func Start(handler http.Handler) {
	svr := &http.Server{
		Addr:    ":3001",
		Handler: handler,
	}

	go func() {
		svr.ListenAndServe()
	}()

	defer stop(svr)

	log.Printf("Started server on %s", svr.Addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")
}

func stop(svr *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
