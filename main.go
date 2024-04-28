package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Detmon0410/assessment-tax/Route"
)

func main() {

	// Accessing environment variables
	port := os.Getenv("PORT")

	// Creating routes
	echo := Route.GetRoutes()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: echo,
	}

	go func() {
		fmt.Printf("Server is running on port %s...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not start server: %v\n", err)
			os.Exit(1)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Shutting down the server...")

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Could not shutdown server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server shutdown completed.")
}
