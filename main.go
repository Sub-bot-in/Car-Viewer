package main

import (
	"cars/internal/handlers"
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
