package pagination

import (
	"math"
	"net/http"
	"sort"
	"strconv"

	"cars/catalog"
)

type PageData struct {
	Cars          []catalog.Model
	Manufacturers []catalog.Manufacturers
	Categories    []catalog.Categories

	Page  int
	Pages []PageLink

	SelectedManufacturer string
	SelectedCategory     string
	SelectedYear         string
	Years                []int

	CompareIDs     map[int]bool
	CompareCount   int
	MaxCompareCars int
}

type PageLink struct {
	Number int
	IsDots bool
}

func buildYears(cars []catalog.Model) []int {
	seen := make(map[int]bool) //to hold unique years
	var years []int

	for _, car := range cars {
		if !seen[car.Year] {
			seen[car.Year] = true
			years = append(years, car.Year)
		}
	}

	sort.Ints(years)

	return years
}

func buildPages(currentPage, totalPages int) []PageLink {
	const sidePages = 2

	var pages []PageLink

	addPage := func(n int) {
		pages = append(pages, PageLink{
			Number: n,
			IsDots: false,
		})
	}

	addDots := func() {
		pages = append(pages, PageLink{
			IsDots: true,
		})
	}

	if totalPages <= 1 {
		addPage(1)
		return pages
	}

	// Если страниц мало, показываем все
	if totalPages <= 7 {
		for i := 1; i <= totalPages; i++ {
			addPage(i)
		}
		return pages
	}

	addPage(1)

	start := currentPage - sidePages
	end := currentPage + sidePages

	if start < 2 {
		start = 2
	}

	if end > totalPages-1 {
		end = totalPages - 1
	}

	if start > 2 {
		addDots()
	}

	for i := start; i <= end; i++ {
		addPage(i)
	}

	if end < totalPages-1 {
		addDots()
	}

	addPage(totalPages)

	return pages
}

func PaginateCars(data *PageData, r *http.Request) {
	const perPage = 20

	data.Years = buildYears(data.Cars)

	page := getCurrentPage(r)

	manufacturer := r.URL.Query().Get("manufacturer")
	category := r.URL.Query().Get("category")
	yearStr := r.URL.Query().Get("year")

	data.SelectedManufacturer = manufacturer
	data.SelectedCategory = category
	data.SelectedYear = yearStr

	filtered := make([]catalog.Model, 0, len(data.Cars))

	for _, car := range data.Cars {

		if manufacturer != "" && car.ManufacturerID != findManufacturerIDByName(data, manufacturer) {
			continue
		}

		if category != "" && car.CategoryID != findCategoryIDByName(data, category) {
			continue
		}

		if yearStr != "" && strconv.Itoa(car.Year) != yearStr {
			continue
		}

		filtered = append(filtered, car)
	}

	totalCars := len(filtered)
	totalPages := calculateTotalPages(totalCars, perPage)

	if page > totalPages {
		page = totalPages
	}

	start, end := getPageBounds(page, perPage, totalCars)

	data.Cars = filtered[start:end]
	data.Page = page
	data.Pages = buildPages(page, totalPages)
}

func getCurrentPage(r *http.Request) int {
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}

	return page
}

func calculateTotalPages(totalItems, perPage int) int {
	if totalItems == 0 {
		return 1
	}

	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}

func getPageBounds(page, perPage, totalItems int) (int, int) {
	start := (page - 1) * perPage
	end := start + perPage

	if end > totalItems {
		end = totalItems
	}

	return start, end
}

func findManufacturerIDByName(data *PageData, name string) int {
	for _, m := range data.Manufacturers {
		if m.Name == name {
			return m.ID
		}
	}
	return 0
}

func findCategoryIDByName(data *PageData, name string) int {
	for _, c := range data.Categories {
		if c.Name == name {
			return c.ID
		}
	}
	return 0
}
