package main

import (
	"fmt"
	"os"

	database "github.com/ghostcode-sys/m/v2/Database"
	routing "github.com/ghostcode-sys/m/v2/Routing"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Connection", os.Getenv("DBNAME"), os.Getenv("DB_URL"))
	defer database.CloseConnection()

	_, err := database.GetDatabaseConnection()

	if err != nil {
		fmt.Println("DB connection Failed")
	}

	if os.Getenv("DOCKER") != "yes" {
		err = godotenv.Load()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	r := routing.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	if os.Getenv("DOCKER") != "yes" {
		r.Run(":8082")
	} else {
		r.Run(":8081")
	}
}
