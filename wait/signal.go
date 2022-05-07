package wait

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
var handler = make(chan struct{})

// SetupStopSignal registered for SIGTERM and SIGINT
func SetupStopSignal() (<-chan struct{}, context.Context) {
	close(handler)

	stopCh := make(chan struct{})
	stopCtx, cancel := context.WithCancel(context.Background())

	signCh := make(chan os.Signal, 2)
	signal.Notify(signCh, shutdownSignals...)
	go func() {
		s1 := <-signCh
		log.Printf("Received signal [%v], beginning shutdown process...\n", s1)
		cancel()
		close(stopCh)

		// Exit directly when received second signal
		s2 := <-signCh
		log.Printf("Received signal [%v] again, will be force to exit", s2)
		os.Exit(1)
	}()
	return stopCh, stopCtx
}
