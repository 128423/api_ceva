package models

import (
	"context"

	"github.com/luis300997/api_ceva/database"

	"time"
)

//Temperatura  Model da temperatura
type Temperatura struct {
	Data        time.Time `json:"data"`
	Temperatura float32   `json:"temperatura"`
}

//GetAllTemperatura retorna as ultimas 5
func GetAllTemperatura() (map[string]Temperatura, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	var a map[string]Temperatura
	q := db.NewRef("temp").OrderByKey().LimitToFirst(5)
	if err := q.Get(context.Background(), &a); err != nil {
		return nil, err
	}
	return a, nil
}
