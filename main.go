package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	err := ConnectDB()
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&User{})
	router := gin.Default()

	router.POST("/signup", SignUp)
	router.POST("/login", Login)
	router.GET("/protected", Protected)

	router.Run(":8080")
}
