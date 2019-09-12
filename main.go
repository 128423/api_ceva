package main

import (
	"os"

	"api_ceva/controlers"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"

	"github.com/gin-contrib/cors"
)

func main() {
	gotenv.Load()
	gin.SetMode(os.Getenv("GIN_MODE"))

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "GET", "DELETE", "POST"},
		AllowHeaders:    []string{"Content-type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length", "Content-type"},
		MaxAge:          36000,
	}))

	router.GET("/temperatura", controlers.GetLast5Temps)

	router.Run()
}
