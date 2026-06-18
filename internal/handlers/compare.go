package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cars/catalog"
	"cars/internal/api"
)

const compareCookieName = "compare_cars"
const maxCompareCars = 4

type ComparePageData struct {
	Cars          []catalog.Model
	Manufacturers map[int]catalog.Manufacturers
	Categories    map[int]catalog.Categories
}

func AddToCompareHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/compare/add/")

		ids := getCompareIDs(r)

		if !containsID(ids, id) && len(ids) < maxCompareCars {
			ids = append(ids, id)
		}

		setCompareCookie(w, ids)

		redirectURL := r.Referer()
		if redirectURL == "" {
			redirectURL = "/compare"
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func RemoveFromCompareHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/compare/remove/")

		ids := getCompareIDs(r)
		var updated []string

		for _, existingID := range ids {
			if existingID != id {
				updated = append(updated, existingID)
			}
		}

		setCompareCookie(w, updated)

		http.Redirect(w, r, "/compare", http.StatusSeeOther)
	}
}

func ClearCompareHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCompareCookie(w, []string{})
		http.Redirect(w, r, "/compare", http.StatusSeeOther)
	}
}

func CompareHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := getCompareIDs(r)

		data := api.LoadCarsAPI()

		var selectedCars []catalog.Model

		for _, id := range ids {
			for _, car := range data.Cars {
				if fmt.Sprint(car.ID) == id {
					selectedCars = append(selectedCars, car)
					break
				}
			}
		}

		manufacturersMap := make(map[int]catalog.Manufacturers)
		for _, manufacturer := range data.Manufacturers {
			manufacturersMap[manufacturer.ID] = manufacturer
		}

		categoriesMap := make(map[int]catalog.Categories)
		for _, category := range data.Categories {
			categoriesMap[category.ID] = category
		}

		pageData := ComparePageData{
			Cars:          selectedCars,
			Manufacturers: manufacturersMap,
			Categories:    categoriesMap,
		}

		err := tmpl.ExecuteTemplate(w, "compare.html", pageData)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getCompareIDs(r *http.Request) []string {
	cookie, err := r.Cookie(compareCookieName)
	if err != nil || cookie.Value == "" {
		return []string{}
	}

	ids := strings.Split(cookie.Value, ",")

	var result []string
	for _, id := range ids {
		id = strings.TrimSpace(id)

		if id == "" {
			continue
		}

		if _, err := strconv.Atoi(id); err == nil {
			result = append(result, id)
		}
	}

	return result
}

func setCompareCookie(w http.ResponseWriter, ids []string) {
	value := strings.Join(ids, ",")

	http.SetCookie(w, &http.Cookie{
		Name:     compareCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func containsID(ids []string, id string) bool {
	for _, existingID := range ids {
		if existingID == id {
			return true
		}
	}

	return false
}
