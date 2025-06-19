package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v3"
)

var (
	reqInfoColl *mongo.Collection
	reqDataColl *mongo.Collection
)

func InitCollections(client *mongo.Client) {
	db := client.Database("postmanData")
	reqInfoColl = db.Collection("reqInfo")
	reqDataColl = db.Collection("reqData")
}

func WithCORS(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Accept")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		if handler != nil {
			handler(c)
		}
	}
}

func GetReqInfoWithData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "reqData"},
			{Key: "localField", Value: "requestDataID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "reqDataDetails"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$reqDataDetails"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$requestUrl"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "allRequests", Value: bson.D{
				{Key: "$push", Value: bson.D{
					{Key: "url", Value: "$requestUrl"},
					{Key: "type", Value: "$requestType"},
					{Key: "status", Value: "$requestStatus"},
					{Key: "timestamp", Value: "$timestamp"},
					{Key: "data", Value: "$reqDataDetails"},
				}},
			}},
			{Key: "createdAt", Value: bson.D{{Key: "$max", Value: "$timestamp"}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "endpoint", Value: "$_id"},
			{Key: "_id", Value: 0},
			{Key: "count", Value: 1},
			{Key: "allRequests", Value: 1},
			{Key: "createdAt", Value: 1},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},
	}

	cursor, err := reqInfoColl.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Aggregation error: " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decoding error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func ConvertHandler(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	yamlData, err := yaml.Marshal(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "YAML conversion error"})
		return
	}

	c.Data(http.StatusOK, "text/plain", yamlData)
}
