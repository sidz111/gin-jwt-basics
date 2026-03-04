package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("supersecretkey")

// ---------------- SIGNUP ----------------

func SignUp(c *gin.Context) {

	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hashedPassword)
	user.Created_at = time.Now()
	user.Updated_at = time.Now()

	DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

// ---------------- LOGIN ----------------

func Login(c *gin.Context) {

	var user User
	var foundUser User

	c.BindJSON(&user)

	DB.Where("username = ?", user.Username).First(&foundUser)

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// 🔥 Create JWT Token
	claims := jwt.MapClaims{
		"username": foundUser.Username,
		"user_id":  foundUser.User_id,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(SECRET_KEY)

	foundUser.Token = signedToken
	DB.Save(&foundUser)

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}

// ---------------- PROTECTED ----------------

func Protected(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token"})
		return
	}

	tokenString := strings.Split(authHeader, " ")[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You are authorized"})
}
