package main

import (
	"cars/catalog"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {

	var data catalog.CarsAPI = loadCarsAPI()
	fmt.Println(len(data.Models))
	fmt.Println(len(data.Categories))
	fmt.Println(len(data.Manufacturers))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/cars", carsHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var tmpl *template.Template = template.Must(template.ParseFiles("templates/index.html"))

type PageData struct {
	Cars []catalog.Model
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := loadCarsAPI()

	pageData := PageData{
		Cars: data.Models,
	}

	tmpl.ExecuteTemplate(w, "index.html", pageData)

}

func carsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cars"))
}

func loadCarsAPI() catalog.CarsAPI {
	return catalog.CarsAPI{
		Models:        fetchModels(),
		Categories:    fetchCategories(),
		Manufacturers: fetchManufacturers(),
	}
}
