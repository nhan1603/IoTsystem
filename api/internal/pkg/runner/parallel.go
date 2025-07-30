package runner

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// ExecParallel executes multiple services in parallel bound by a shared context.
// This returns the first non-nil error (if any) from them. If any of the components encounter an error, it will trigger
// shutdown for all the other components.
func ExecParallel(ctx context.Context, services ...func(context.Context) error) {
	log.Println("Starting all services...")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(services)) // if num of goroutines is fixed, should just use Add once.
	for i := range services {
		s := services[i] // https://golang.org/doc/faq#closures_and_goroutines
		go func() {
			defer wg.Done()

			err := s(ctx) // run and wait

			if err != nil {
				log.Printf("service encountered err, %v", err)

				cancel() // if one svc encounters error, terminate everything else
			}
		}()
	}

	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)

	select { // blocks until either os cancel signal is sent or the running context is cancelled.
	case s := <-ch:
		log.Printf("Received signal: [%s] to terminate. Shutting down all services...", s.String())
		cancel()
	case <-ctx.Done():
		log.Println("Context cancelled. Shutting down all services...")
	}

	wg.Wait()

	log.Println("All services shut down")
}
