package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims, err := VerifyToken(tokenString)
		log.Println(claims)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		userID := claims["user_id"].(string)
		c.Set("user_id", userID)
	}
}