package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/Suke2004/students-api/internal/config"
	"github.com/Suke2004/students-api/internal/http/handlers/student"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup

	//setup router
	router := http.NewServeMux() //response w,     request r
	router.HandleFunc("POST /api/students", student.New())

	//setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr))
	fmt.Printf("server started %s\n", cfg.HTTPServer.Addr)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server") // go run ...../main.go -config config.local.yaml    to run the server
		}
	}()
	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //If a response is send when a request being process it wont accept and when we shutdown it wil wait to 5 sec for on going request to proces
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown sucessfully")

}

//Alternative code for above
// package main

// import (
// 	"context"
// 	"fmt"
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/Suke2004/students-api/internal/config"
// )

// func main() {
// 	// Load config
// 	cfg := config.MustLoad()

// 	// Setup router
// 	router := http.NewServeMux() // response w, request r
// 	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Welcome to students api"))
// 	})

// 	// Setup server
// 	server := http.Server{
// 		Addr:    cfg.HTTPServer.Addr,
// 		Handler: router,
// 	}

// 	// Log the server start
// 	slog.Info("Server started", slog.String("address", cfg.HTTPServer.Addr))
// 	fmt.Printf("Server started on %s\n", cfg.HTTPServer.Addr)

// 	// Channel to listen for interrupt signals
// 	done := make(chan os.Signal, 1)
// 	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// 	// Run the server in a goroutine
// 	go func() {
// 		err := server.ListenAndServe()
// 		if err != nil && err != http.ErrServerClosed {
// 			// If the server fails, log the error
// 			slog.Error("Server failed to start", slog.String("error", err.Error()))
// 		}
// 	}()

// 	// Wait for shutdown signal
// 	sigReceived := <-done
// 	slog.Info("Shutdown signal received", slog.String("signal", sigReceived.String()))

// 	// Graceful shutdown with a 5-second timeout
// 	slog.Info("Shutting down the server gracefully")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Gracefully wait for 5 seconds for ongoing requests to complete
// 	defer cancel()

// 	// Attempt to shut down the server
// 	if err := server.Shutdown(ctx); err != nil {
// 		slog.Error("Failed to shut down the server", slog.String("error", err.Error()))
// 	} else {
// 		slog.Info("Server shutdown successfully")
// 	}
// }
