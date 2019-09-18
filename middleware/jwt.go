package middleware

import (
	database "controle-api/Database"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Jwt ...
func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		// flagNewEmpresa := c.GetHeader("flag")

		// // flag, _ := strconv.ParseBool(flagNewEmpresa)

		// // if !flag {
		// // 	return
		// // }

		if tokenString != "" {
			bearerToken := strings.Split(tokenString, " ")

			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"authErrors": []string{"Invalid authorization token"}})
					c.Abort()
					return
				}

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					usu := claims["sub"].(float64)
					emp := claims["emp"].(float64)
					cpfCnpj := claims["cpf_cnpj"].(float64)
					remember := claims["remember"].(bool)

					if claims["sub"] == nil || claims["emp"] == nil || claims["cpf_cnpj"] == nil {
						c.JSON(http.StatusBadRequest, gin.H{"authErrors": []string{"Invalid payload"}})
						c.Abort()
						return
					}

					db, er := database.SetupDB()
					if er != nil {
						c.JSON(400, gin.H{"authErrors": []string{"Error on database connection"}})
						c.Abort()
						return
					}
					defer db.Close()

					var t string
					db = db.Table("controle.con_tokens").
						Where("tok_token = ?", bearerToken[1]).
						Where("tok_data_exc is null")

					row := db.Select("tok_token").Row()
					row.Scan(&t)
					if len(t) == 0 {
						c.JSON(400, gin.H{"authErrors": []string{"Invalid authorization token"}})
						c.Abort()
						return
					}

					c.Set("usu_codigo", int(usu))
					c.Set("emp_codigo", int(emp))
					c.Set("cpf_cnpj", int(cpfCnpj))
					c.Set("remember", remember)
					return
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"authErrors": []string{err.Error()}})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"authErrors": []string{"Invalid authorization token"}})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"authErrors": []string{"An authorization header is required"}})
			c.Abort()
			return
		}
	}
}
