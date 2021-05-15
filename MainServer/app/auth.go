package app

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-contacts/models"
	"net/http"
	"os"
	"strings"
)

func JwtAuthentication(c *gin.Context)  {
	tokenHeader := c.GetHeader("Authorization") //Grab the token from the header

	if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
		c.Status(http.StatusForbidden)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, map[string]string{"message":"Missing auth token"})
		c.Abort()
		return
	}

	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		c.Status(http.StatusForbidden)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, map[string]string{"message":"Invalid/Malformed auth token"})
		c.Abort()
		return
	}

	tokenPart := splitted[1] //Grab the token part, what we are truly interested in
	tk := &models.Token{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	if err != nil { //Malformed token, returns with http code 403 as usual
		c.Status(http.StatusForbidden)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, "Malformed authentication token")
		c.Abort()
		return
	}

	if !token.Valid { //Token is invalid, maybe not signed on this server
		c.Status(http.StatusForbidden)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, "Token is not valid.")
		c.Abort()
		return
	}

	c.Set("user", tk.UserId)
	c.Next()
}
