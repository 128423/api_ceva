package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//Google middleware para altenticar usuarios com oauth
func Google() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
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
					id := claims["id"].(string)
					email := claims["email"].(string)
					given_name := claims["given_name"].(string)
					family_name := claims["family_name"].(string)
					Picture := claims["picture"].(string)
					Locale := claims["locale"].(string)
					Hd := claims["hd"].(string)
					tokenGoogle := claims["tokenGoogle"].(string)

					c.Set("id", id)
					c.Set("token", bearerToken[1])
					c.Set("email", email)
					c.Set("given_name", given_name)
					c.Set("family_name", family_name)
					c.Set("picture", Picture)
					c.Set("locale", Locale)
					c.Set("hd", Hd)
					c.Set("tokenGoogle", tokenGoogle)
					return
				}
				c.JSON(http.StatusUnauthorized, gin.H{"authErrors": []string{err.Error()}})
				c.Abort()
				return

			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"authErrors": []string{"Invalid payload"}})
		c.Abort()
		return
	}

}
