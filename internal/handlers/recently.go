package handlers

import (
	"cars/catalog"
	"net/http"
	"strconv"
	"strings"
)

const recentlyViewedCookieName = "recently_viewed"
const maxRecentlyViewedStored = 6
const maxRecentlyViewedShown = 3

func GetRecentlyViewedIDs(r *http.Request) (res []string) {
	cookie, err := r.Cookie(recentlyViewedCookieName)
	if err != nil || cookie.Value == "" {
		return []string{}
	}

	cookieIds := strings.Split(cookie.Value, ",")

	for _, id := range cookieIds {
		id = strings.TrimSpace(id)

		if id != "" {
			res = append(res, id)
		}
	}

	return res
}

func AddRecentlyViewed(w http.ResponseWriter, r *http.Request, id string) {
	cookieIds := GetRecentlyViewedIDs(r)

	updated := []string{id}

	for _, existingID := range cookieIds {
		if existingID == id {
			continue
		}

		updated = append(updated, existingID)

		if len(updated) == maxRecentlyViewedStored {
			break
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     recentlyViewedCookieName,
		Value:    strings.Join(updated, ","),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
func GetRecentlyViewedCars(cars []catalog.Model, cookieIds []string, currentCarID int) (res []catalog.Model) {
	for _, id := range cookieIds {
		for _, car := range cars {
			if car.ID == currentCarID {
				continue
			}

			if strconv.Itoa(car.ID) == id {
				res = append(res, car)
				break
			}
		}

		if len(res) == maxRecentlyViewedShown {
			break
		}
	}

	return res
}
