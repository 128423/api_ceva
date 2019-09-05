package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/luis300997/api_ceva/controlers"

	"github.com/gin-contrib/cors"
	"github.com/subosito/gotenv"
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
