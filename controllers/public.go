package controllers

import (
	"go-jwt/auth"
	"go-jwt/database"
	"go-jwt/models"
	"log"
	"github.com/gin-gonic/gin"
)

type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}

func Signup(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		log.Printf("Error while binding JSON: %v", err)
		c.JSON(400, gin.H{
			"error": "Invalid request",
		})
		c.Abort()
		return
	}
	err = user.HashPassword(user.Password)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(500, gin.H{
			"error": "Internal creating error",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"message": "User created successfully",
	})
}
func Login(c *gin.Context) {
	var payload LoginPayload
	var user models.User

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request",
		})
		c.Abort()
		return
	}
	result := database.GlobalDB.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": "Invalid email or password",
		})
		c.Abort()
		return
	}

	err = user.CheckPassword(payload.Password)
	if err != nil {
		c.JSON(404, gin.H{
			"Error": "Invalid email or password",
		})
		c.Abort()
		return
	}
	jwtWrapper := auth.JwtWrapper{
		SecretKey:         "secretKey",
		Issuer:            "AuthService",
		ExpirationMinutes: 1,
		ExpirationHours:   12,
	}
	signedToken, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error while signing token",
		})
		c.Abort()
		return
	}
	signedToken, err1 := jwtWrapper.RefreshToken(user.Email)
	if err1 != nil {
		log.Println("Error while refreshing token: %v", err)
		c.JSON(500, gin.H{
			"error": "Error while signing token",
		})
		c.Abort()
		return
	}
	tokenResponse := LoginResponse{
		Token:        signedToken,
		RefreshToken: signedToken,
	}
	c.JSON(200, tokenResponse)
}
