package main

import (
	"cars/catalog"
	"math/rand/v2"
)

func getRelatedCars(cars []catalog.Model, currentCar catalog.Model, limit int) []catalog.Model {
	var related []catalog.Model

	for _, car := range cars {
		if car.ID == currentCar.ID {
			continue
		}

		if car.ManufacturerID == currentCar.ManufacturerID {
			related = append(related, car)
		}
	}

	rand.Shuffle(len(related), func(i, j int) {
		related[i], related[j] = related[j], related[i]
	})

	if len(related) > limit {
		related = related[:limit]
	}

	return related
}
