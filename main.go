package main

import (
	"cars/catalog"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/cars", carsHandler)
	http.HandleFunc("/spec/", specHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var tmpl *template.Template = template.Must(template.ParseGlob("templates/*.html"))

type PageData struct {
	Cars          []catalog.Model
	Manufacturers []catalog.Manufacturers
	Categories    []catalog.Categories
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := loadCarsAPI()

	pageData := PageData{
		Cars:          data.Models,
		Manufacturers: data.Manufacturers,
		Categories:    data.Categories,
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

type CarWithDetails struct {
	Car          catalog.Model
	Manufacturer catalog.Manufacturers
	Categorie    catalog.Categories
}

func specHandler(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/spec/")

	data := loadCarsAPI()

	var car catalog.Model
	var manufacturer catalog.Manufacturers
	var categorie catalog.Categories

	for _, c := range data.Models {
		if fmt.Sprint(c.ID) == id {
			car = c
			break
		}
	}

	for _, c := range data.Manufacturers {
		if fmt.Sprint(c.ID) == id {
			manufacturer = c
			break
		}
	}

	for _, c := range data.Categories {
		if fmt.Sprint(c.ID) == id {
			categorie = c
			break
		}
	}

	specPage := CarWithDetails{
		Car:          car,
		Manufacturer: manufacturer,
		Categorie:    categorie,
	}
	err := tmpl.ExecuteTemplate(w, "car.html", specPage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
