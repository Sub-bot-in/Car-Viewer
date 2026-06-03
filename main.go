package main

import (
	//"bytes"
	//m "cars/models"
	//"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	//"strings"
)

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/cars", dataHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API cars http://localhost:8080/api/cars"))
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	info, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(info))

	// 	var data m.CarsAPI

	// 	json.Unmarshal(info, &data)

	// 	// err = json.NewDecoder(resp.Body).Decode(&data)
	// 	// if err != nil {
	// 	// 	http.Error(w, "bad json", http.StatusInternalServerError)
	// 	// 	return
	// 	// }

	// 	//w.Header().Set("Content-Type", "application/json")
	// 	//json.NewEncoder(w).Encode(data)

	// 	data1, _ := json.Marshal(data)

	// 	resp2, err := http.Post(
	// 		"http://localhost:8080/api/cars",
	// 		"application/json",
	// 		bytes.NewBuffer(data1),
	// 	)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer resp2.Body.Close()

	//	for _, car := range data1 {
	//		go func(c m.CarsAPI) {
	//			sendToServer(c)
	//		}(data)
	//	}
}
