package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gonzalezryan03/gopubsub/internal/pkg/httpserver"
	"github.com/gonzalezryan03/gopubsub/internal/pubsub"
)

func main() {
	// 1. Create the pubsub service
	psService := pubsub.NewService()

	// 2. Create the http server
	httpServer := httpserver.NewServer(":8080", psService)

	// 3. Set up signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server in separate goroutine
	go func() {
		if err := httpServer.Start(ctx); err != nil {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}()

	fmt.Println("GoPubSub server running on port 8080")

	// 4. Wait for signal to stop
	<-ctx.Done()

	// 5. Shutdown the server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server failed to shutdown: %v", err)
	}

	fmt.Println("GoPubSub server shutdown complete")

}
