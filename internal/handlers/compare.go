package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"cars/catalog"
	"cars/internal/api"
)

const compareCookieName = "compare_cars"
const MaxCompareCars = 4

type ComparePageData struct {
	Cars          []catalog.Model
	Manufacturers map[int]catalog.Manufacturers
	Categories    map[int]catalog.Categories
}

func AddToCompareHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/compare/add/")

		ids := GetCompareIDs(r)

		if !containsID(ids, id) && len(ids) < MaxCompareCars {
			ids = append(ids, id)
		}

		setCompareCookie(w, ids)

		query := url.Values{}

		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		query.Set("page", page)

		manufacturer := r.URL.Query().Get("manufacturer")
		category := r.URL.Query().Get("category")
		year := r.URL.Query().Get("year")

		if manufacturer != "" {
			query.Set("manufacturer", manufacturer)
		}

		if category != "" {
			query.Set("category", category)
		}

		if year != "" {
			query.Set("year", year)
		}

		from := r.URL.Query().Get("from")

		var redirectURL string

		if from == "details" {
			redirectURL = "/specifications/" + id + "?" + query.Encode()
		} else {
			redirectURL = "/?" + query.Encode() + "#car-" + id
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func RemoveFromCompareHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/compare/remove/")

		ids := GetCompareIDs(r)

		var updated []string

		for _, existingID := range ids {
			if existingID != id {
				updated = append(updated, existingID)
			}
		}

		setCompareCookie(w, updated)

		query := url.Values{}

		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}

		query.Set("page", page)

		manufacturer := r.URL.Query().Get("manufacturer")
		category := r.URL.Query().Get("category")
		year := r.URL.Query().Get("year")

		if manufacturer != "" {
			query.Set("manufacturer", manufacturer)
		}

		if category != "" {
			query.Set("category", category)
		}

		if year != "" {
			query.Set("year", year)
		}

		from := r.URL.Query().Get("from")

		var redirectURL string

		if from == "details" {
			redirectURL = "/specifications/" + id + "?" + query.Encode()
		} else if from == "catalog" {
			redirectURL = "/?" + query.Encode() + "#car-" + id
		} else {
			redirectURL = "/compare"
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func ClearCompareHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCompareCookie(w, []string{})

		redirectURL := r.URL.Query().Get("return")
		if redirectURL == "" {
			redirectURL = "/compare"
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func CompareHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := GetCompareIDs(r)

		data, err := api.LoadCarsAPI()

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			tmpl.ExecuteTemplate(w, "maintenance.html", nil)
			return
		}

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

		err = tmpl.ExecuteTemplate(w, "compare.html", pageData)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetCompareIDs(r *http.Request) []string {
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
	if slices.Contains(ids, id) {
		return true
	}
	return false
}

func GetCompareIDMap(r *http.Request) map[int]bool {
	ids := GetCompareIDs(r)

	result := make(map[int]bool)

	for _, id := range ids {
		numID, err := strconv.Atoi(id)
		if err != nil {
			continue
		}

		result[numID] = true
	}

	return result
}
