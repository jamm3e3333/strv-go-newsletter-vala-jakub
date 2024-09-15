package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("starting...")
	notifyShutdownChan := make(chan os.Signal, 1)
	signal.Notify(notifyShutdownChan, []os.Signal{os.Interrupt, os.Kill}...)
	go func() {
		<-notifyShutdownChan
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Shutting down...")
}
