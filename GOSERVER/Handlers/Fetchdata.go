package handlers

import (
	"context"
	"log"
	"os"
	"time"

	database "github.com/ghostcode-sys/m/v2/Database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchData(c *gin.Context) {

	inputParams := ExtractParams(c)

	dbName := os.Getenv("DBNAME")
	testData := os.Getenv("TESTDATA")

	client, err := database.GetDatabaseConnection()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	pattern := inputParams["search"]
	regex := bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}

	cursor, queryErr := client.Database(dbName).Collection(testData).Find(ctx, bson.M{"api_group": regex})

	if queryErr != nil {
		c.JSON(500, gin.H{
			"error": queryErr.Error(),
		})
		return
	}

	queryResult := make([]primitive.M, 0)

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		queryResult = append(queryResult, result)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, queryResult)

}
