package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	database "github.com/ghostcode-sys/m/v2/Database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WriteFile(c *gin.Context) {

	inputParams := ExtractParams(c)
	fmt.Println(inputParams)

	output := map[string]string{
		"status":  "200",
		"message": "success",
	}
	client, err := database.GetDatabaseConnection()
	if err != nil {
		output["message"] = "Failure1: " + err.Error()
		output["status"] = "500"
		c.JSON(http.StatusInternalServerError, output)
		return
	}
	dbName := os.Getenv("DBNAME")
	reqInfo := os.Getenv("REQINFO")
	reqData := os.Getenv("REQDATA")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	var RequestData database.ReqData

	json_err := json.Unmarshal([]byte(inputParams["params"]), &RequestData.Params)
	if json_err != nil {
		output["message"] = "Failure2: " + json_err.Error()
		output["status"] = "500"
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	json_err = json.Unmarshal([]byte(inputParams["response"]), &RequestData.Response)

	fmt.Println(RequestData)
	result, queryErr := client.Database(dbName).Collection(reqData).InsertOne(ctx, RequestData)
	if queryErr != nil {
		output["message"] = "Failure3: " + queryErr.Error()
		output["status"] = "500"
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	objectID := result.InsertedID.(primitive.ObjectID)
	output["reqDataId"] = objectID.String()

	var RequestInfo database.ReqInfo

	RequestInfo.RequestDataID = objectID
	RequestInfo.RequestUrl = inputParams["url"]
	RequestInfo.Timestamp = time.Now()
	RequestInfo.RequestType = inputParams["type"]
	RequestInfo.RequestStatus, err = strconv.Atoi(inputParams["status"])
	RequestInfo.ResponseTime, err = strconv.Atoi(inputParams["responsetime"])
	if err != nil {
		output["message"] = "Failure4: " + err.Error()
		output["status"] = "500"
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel2()

	result, queryErr = client.Database(dbName).Collection(reqInfo).InsertOne(ctx2, RequestInfo)
	if queryErr != nil {
		output["message"] = "Failure5: " + queryErr.Error()
		output["status"] = "500"
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	objectID2 := result.InsertedID.(primitive.ObjectID)

	output["reqInfoId"] = objectID2.String()

	c.JSON(http.StatusOK, output)
}
