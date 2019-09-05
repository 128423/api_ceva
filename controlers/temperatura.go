package controlers

import (
	"github.com/gin-gonic/gin"
	"github.com/luis300997/api_ceva/models"
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
