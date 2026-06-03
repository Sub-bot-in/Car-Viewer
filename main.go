package main

import (
	"cars/catalog"
	"fmt"
	"log"
	"net/http"
)

func main() {

	var data catalog.CarsAPI = loadCarsAPI()
	fmt.Println(data)

	http.HandleFunc("/", homeHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API cars"))
}

func loadCarsAPI() catalog.CarsAPI {
	return catalog.CarsAPI{
		Models:        fetchModels(),
		Categories:    fetchCategories(),
		Manufacturers: fetchManufacturers(),
	}
}
