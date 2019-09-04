package database

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
)

//Connect Coneção do banco
func Connect() (*db.Client, error) {
	ctx := context.Background()
	config := &firebase.Config{
		DatabaseURL: "https://cervejaproject.firebaseio.com",
	}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, nil
	}
	client, err := app.Database(ctx)

	if err != nil {
		return nil, nil
	}
	return client, nil
}
