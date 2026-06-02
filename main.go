package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/cars", dataHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API carshttp://localhost:8080/api/cars"))
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("api/data.json")
	if err != nil {
		http.Error(w, "file not found", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var data CarsAPI

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		http.Error(w, "bad json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
