package controlers

import (
	"api_ceva/models"

	"github.com/gin-gonic/gin"
)

//GetLast5Temps handler
func GetLast5Temps(c *gin.Context) {
	ret, err := models.GetAllTemperatura()
	if err != nil {
		c.JSON(500, gin.H{"errors": []string{"erro ao pegar dados " + err.Error()}})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"data": ret})
	return
}
