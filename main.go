package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/low4ey/OJ/Golang-backend/routes"
	"github.com/low4ey/OJ/Golang-backend/utils"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	port := os.Getenv("PORT")
	// config, err := config.LoadConfig(".")
	// if err != nil {
	// 	fmt.Println("Environment Variable Failed Loading")
	// 	os.Exit(1)
	// }
	// port := config.PORT
	if port == "" {
		port = "8080"
	}
	fmt.Println("And the Port Is : " + port)
	router := gin.Default()
	router.Use(utils.CORSMiddleware())
	routes.SubmissionRoutes(router)
	log.Fatal(router.Run(":" + port))
}
