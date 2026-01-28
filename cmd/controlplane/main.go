package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhangbin/Velo/internal/xds"
)

func main() {
	// Set up basic logger
	logger := log.New(os.Stdout, "[VELO] ", log.LstdFlags)

	// Create context that cancels on OS signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Printf("Received signal: %v, shutting down...", sig)
		cancel()
	}()

	// Initialize Control Plane
	cp, err := xds.NewControlPlane(ctx, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize control plane: %v", err)
	}

	// Start the Control Plane
	// This blocks until error or context cancellation (if modified to support it, 
	// but currently Run() blocks on Serve). 
	// In a real app we might run this in a goroutine, but for now blocking is fine.
	logger.Println("Starting Velo Control Plane...")
	if err := cp.Run(); err != nil {
		logger.Fatalf("Control plane error: %v", err)
	}
}
