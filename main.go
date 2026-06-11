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
	//	http.HandleFunc("/cars", carsHandler)
	http.HandleFunc("/specifications/", specHandler)
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

	data := PageData{}

	errCh := make(chan error, 3) // буфер важен!
	doneCh := make(chan struct{}, 3)

	go func() {
		err := fetchJSON(baseURL+"/models", &data.Cars)
		errCh <- err
		doneCh <- struct{}{}
	}()

	go func() {
		err := fetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
		errCh <- err
		doneCh <- struct{}{}
	}()

	go func() {
		err := fetchJSON(baseURL+"/categories", &data.Categories)
		errCh <- err
		doneCh <- struct{}{}
	}()

	// 🔥 fan-in (ждём 3 результата)

	for i := 0; i < 3; i++ {
		<-doneCh
	}

	close(errCh)

	// проверяем ошибки
	for err := range errCh {
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// func carsHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("cars"))
// }

func loadCarsAPI() PageData {
	var data PageData

	err := fetchJSON(baseURL+"/models", &data.Cars)
	if err != nil {
		log.Println(err)
	}

	err = fetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
	if err != nil {
		log.Println(err)
	}

	err = fetchJSON(baseURL+"/categories", &data.Categories)
	if err != nil {
		log.Println(err)
	}

	return data
}

type CarWithDetails struct {
	Car          catalog.Model
	Manufacturer catalog.Manufacturers
	Category     catalog.Categories
}

func specHandler(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/specifications/")

	data := loadCarsAPI()

	var car catalog.Model
	var manufacturer catalog.Manufacturers
	var category catalog.Categories

	for _, c := range data.Cars {
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
			category = c
			break
		}
	}

	specPage := CarWithDetails{
		Car:          car,
		Manufacturer: manufacturer,
		Category:     category,
	}
	err := tmpl.ExecuteTemplate(w, "car.html", specPage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
