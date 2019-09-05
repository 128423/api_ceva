package models

import (
	"context"
	"log"

	"github.com/luis300997/api_ceva/database"

	"time"
)

//Temperatura  Model da temperatura
type Temperatura struct {
	Data        time.Time `json:"data"`
	Temperatura float32   `json:"temperatura"`
}

//GetAllTemperatura retorna as ultimas 5
func GetAllTemperatura() ([]Temperatura, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	var temps []Temperatura
	var a map[Temperatura]interface{}
	q := db.NewRef("temp").OrderByKey().LimitToFirst(5)
	if err := q.Get(context.Background(), &a); err != nil {
		return nil, err
	}

	for aux := range a {
		log.Println(aux)
	}

	return temps, nil
}
