package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	database "github.com/ghostcode-sys/m/v2/Database"
	handlers "github.com/ghostcode-sys/m/v2/Handlers"
	routing "github.com/ghostcode-sys/m/v2/Routing"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("DOCKER") != "yes" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	fmt.Println("Connection", os.Getenv("DBNAME"), os.Getenv("DB_URL"))

	client, err := database.GetDatabaseConnection()
	if err != nil {
		log.Fatal("DB connection Failed:", err) 
	}
	defer database.CloseConnection()

	handlers.InitCollections(client)

	r := routing.SetupRouter()

	r.GET("/assets/*filepath", func(c *gin.Context) {
		file := c.Param("filepath")
		fullPath := "./Static/assets" + file

		fmt.Println("Request for asset:", fullPath)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Println("File not found:", fullPath)
			c.Status(404)
			return
		}

		if strings.HasSuffix(file, ".css") {
			c.Header("Content-Type", "text/css")
		} else if strings.HasSuffix(file, ".js") {
			c.Header("Content-Type", "application/javascript")
		}

		c.File(fullPath)
	})

	r.NoRoute(func(c *gin.Context) {
		c.File("./Static/index.html")
	})

	r.GET("/react", func(c *gin.Context) {
		c.File("./Static/index.html")
	})
	r.GET("/old", func(c *gin.Context) {
		c.File("./Frontend/index.html")
	})
	r.LoadHTMLGlob("Frontend/*")

	port := ":8081"
	fmt.Printf("Server running on port %s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
