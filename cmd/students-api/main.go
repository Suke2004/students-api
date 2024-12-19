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
	"github.com/Suke2004/students-api/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initiated", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router
	router := http.NewServeMux() //response w,     request r
	router.HandleFunc("POST /api/students/", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students/", student.GetList(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))
	router.HandleFunc("DELETE /api/students/", student.DeleteAll(storage))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateById(storage))

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


//Mongo db

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/Suke2004/students-api/internal/config"
// 	"github.com/Suke2004/students-api/internal/handler"
// 	"github.com/Suke2004/students-api/internal/storage/mongodb"
// )

// func main() {
// 	// Load configuration
// 	cfg := &config.Config{
// 		Port: "8080",
// 		MongoDB: struct {
// 			URI      string
// 			Database string
// 		}{
// 			URI:      "mongodb://localhost:27017",
// 			Database: "students_api",
// 		},
// 	}

// 	// Initialize MongoDB
// 	db, err := mongodb.New(cfg)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize MongoDB: %v", err)
// 	}
// 	defer func() {
// 		if err := db.Client.Disconnect(nil); err != nil {
// 			log.Fatalf("Failed to disconnect MongoDB client: %v", err)
// 		}
// 	}()

// 	// Initialize API handlers
// 	studentHandler := handler.NewStudentHandler(db)

// 	// Setup routes
// 	http.HandleFunc("/students", studentHandler.HandleStudents)       // GET/POST
// 	http.HandleFunc("/students/", studentHandler.HandleStudentById)  // GET/PUT/DELETE

// 	// Start server
// 	fmt.Printf("Server is running on port %s\n", cfg.Port)
// 	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
// }


//PostgreSQL

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/Suke2004/students-api/internal/config"
// 	"github.com/Suke2004/students-api/internal/handler"
// 	"github.com/Suke2004/students-api/internal/storage/postgresql"
// )

// func main() {
// 	// Load configuration
// 	cfg := &config.Config{
// 		Postgres: struct {
// 			Host     string
// 			Port     string
// 			User     string
// 			Password string
// 			DbName   string
// 			SSLMode  string
// 		}{
// 			Host:     "localhost",
// 			Port:     "5432",
// 			User:     "your_username",
// 			Password: "your_password",
// 			DbName:   "students_api",
// 			SSLMode:  "disable",
// 		},
// 	}

// 	// Initialize PostgreSQL
// 	db, err := postgresql.New(cfg)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
// 	}
// 	defer db.Db.Close()

// 	// Initialize API handlers
// 	studentHandler := handler.NewStudentHandler(db)

// 	// Setup routes
// 	http.HandleFunc("/students", studentHandler.HandleStudents)       // GET/POST
// 	http.HandleFunc("/students/", studentHandler.HandleStudentById)  // GET/PUT/DELETE

// 	// Start server
// 	fmt.Printf("Server is running on port %s\n", "8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
