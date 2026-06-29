package main

import (
	"cars/internal/handlers"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", handlers.HomeHandler(tmpl))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/specifications/", handlers.SpecHandler(tmpl))

	http.HandleFunc("/compare", handlers.CompareHandler(tmpl))
	http.HandleFunc("/compare/add/", handlers.AddToCompareHandler())
	http.HandleFunc("/compare/remove/", handlers.RemoveFromCompareHandler())
	http.HandleFunc("/compare/clear", handlers.ClearCompareHandler())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		fmt.Println("Server started on http://localhost:8080")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)

		if err := srv.Close(); err != nil {
			log.Printf("Server close failed: %v", err)
		}
	}

	fmt.Println("Server stopped.")

}
