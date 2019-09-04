package models

import (
	"context"

	"github.com/luis300997/api_ceva/database"

	"time"
)

type Temperatura struct {
	Data        time.Time `json:"data"`
	Temperatura float32   `json:"temperatura"`
}

func getAllTemperatura() ([]Temperatura, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	var temps []Temperatura
	q := db.NewRef("temp").OrderByKey().LimitToFirst(5)
	if err := q.Get(context.Background(), &temps); err != nil {
		return nil, err
	}
	return temps, nil
}
