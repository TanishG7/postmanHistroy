package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	database "github.com/ghostcode-sys/m/v2/Database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCases(c *gin.Context) {
	inputParams := ExtractParams(c)
	requestParams := make(map[string][]string)
	for key, value := range inputParams {
		if key != "url" {
			requestParams[key] = strings.Split(value, ",")
			if key != "method" {
				requestParams[key] = append(requestParams[key], "")
			}
		}
	}

	urlArr := strings.Split(inputParams["url"], ",")

	goUrl, phpUrl := "", ""

	for _, url := range urlArr {
		if strings.Contains(url, "/go/") {
			goUrl = strings.Trim(url, " ")
		} else {
			phpUrl = strings.Trim(url, " ")
		}
	}

	premutatedRequestParams := generatePermutations(requestParams)

	CompleteResponse := make([]map[string]interface{}, 0)

	permutationCount := len(premutatedRequestParams)
	timeinSecond := fmt.Sprintf("%d seconds", permutationCount*2)

	newUUID, _ := exec.Command("uuidgen").Output()
	newUUIDString := string(newUUID)

	c.JSON(200, gin.H{
		"Counts":                  permutationCount,
		"EstimatedTimeTOComplete": timeinSecond,
		"id":                      newUUIDString,
	})
	go func() {
		for idx, requestParam := range premutatedRequestParams {
			fmt.Println(requestParam)
			var goStatus, phpStatus int
			var goResult, phpResult any
			var goErr, phpErr error

			if requestParam["method"] == "" {
				continue
			}

			if strings.ToUpper(requestParam["method"]) == "GET" {
				delete(requestParam, "method")
				if goUrl != "" {
					newGoUrl := strings.Replace(goUrl, "localhost", "host.docker.internal", 1)
					goStatus, goResult, goErr = hitGetRequest(newGoUrl, requestParam)
				}
				if phpUrl != "" {
					newPhpUrl := strings.Replace(phpUrl, "localhost", "host.docker.internal", 1)
					phpStatus, phpResult, phpErr = hitGetRequest(newPhpUrl, requestParam)
				}
			} else if strings.ToUpper(requestParam["method"]) == "POST" {
				delete(requestParam, "method")

				if goUrl != "" {
					newGoUrl := strings.Replace(goUrl, "localhost", "host.docker.internal", 1)
					goStatus, goResult, goErr = hitPostRequest(newGoUrl, requestParam)
				}
				if phpUrl != "" {
					newPhpUrl := strings.Replace(phpUrl, "localhost", "host.docker.internal", 1)
					phpStatus, phpResult, phpErr = hitPostRequest(newPhpUrl, requestParam)
				}
			}
			goErrString, phpErrString := "", ""
			if goErr != nil {
				goErrString = goErr.Error()
			}
			if phpErr != nil {
				phpErrString = phpErr.Error()
			}

			fmt.Println(goErrString, goResult, goStatus)
			fmt.Println(phpErrString, phpResult, phpStatus)
			DataToStore := map[string]interface{}{
				"goUrl":     goUrl,
				"phpUrl":    phpUrl,
				"goErr":     goErrString,
				"phpErr":    phpErrString,
				"goStatus":  goStatus,
				"phpStatus": phpStatus,
				"goResult":  goResult,
				"phpResult": phpResult,
				"params":    requestParam,
				"enteredOn": time.Now(),
				"api_group": newUUIDString,
			}
			fmt.Println("Case Running: ", idx)
			res, storeErr := StoreResult(DataToStore)
			fmt.Println("Result:", res, storeErr)
			fmt.Println("Case stored: ", idx)
			CompleteResponse = append(CompleteResponse, DataToStore)
			time.Sleep(1 * time.Second)

		}
	}()
}

func generatePermutations(input map[string][]string) []map[string]string {
	var keys []string
	var values [][]string

	// Extract keys and corresponding slice values
	for k, v := range input {
		keys = append(keys, k)
		values = append(values, v)
	}

	var result []map[string]string
	temp := make(map[string]string)

	// Recursive function to generate permutations
	var permute func(index int)
	permute = func(index int) {
		if index == len(keys) {
			// Copy the current permutation to result
			perm := make(map[string]string)
			for k, v := range temp {
				perm[k] = v
			}
			result = append(result, perm)
			return
		}

		// Iterate over values of the current key
		for _, val := range values[index] {
			temp[keys[index]] = val
			permute(index + 1)
		}
	}

	permute(0)
	return result
}

func generateCombinations(input map[string][]string) []map[string]string {
	result := make([]map[string]string, 0)
	var keys []string
	var values [][]string

	// Extract keys and corresponding slice values
	for k, v := range input {
		keys = append(keys, k)
		values = append(values, v)
	}
	var combine func(index int, current map[string]string)
	combine = func(index int, current map[string]string) {
		if index == len(keys) {
			// Make a copy to store the current combination
			comb := make(map[string]string)
			for k, v := range current {
				comb[k] = v
			}
			result = append(result, comb)
			return
		}

		// Exclude the current key (optional)
		combine(index+1, current)

		// Include the current key with its possible values
		for _, val := range values[index] {
			current[keys[index]] = val
			combine(index+1, current)
			delete(current, keys[index]) // Backtrack
		}
	}

	combine(0, make(map[string]string))

	return result
}

func hitGetRequest(url string, params map[string]string) (int, any, error) {

	var data any

	var clientTimeout time.Duration

	clientTimeout = 6 * time.Second

	client := &http.Client{
		Timeout: clientTimeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, data, err
	}

	// Add query parameters to the request
	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return 0, data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, string(body), err
	}

	return resp.StatusCode, data, nil
}

func hitPostRequest(url string, params map[string]string) (int, any, error) {
	defer PanicHandler(true, false, "Panic error encountered in ApiPostPlainTextQueryFunction!")
	jsonData, err := json.Marshal(params)
	var returnData interface{}
	if err != nil {
		return 0, returnData, err
	}
	responseStr := ""
	timeout := 3 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return 0, returnData, err
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return 0, returnData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, returnData, err
	}

	responseStr = string(body)

	err = json.Unmarshal(body, &returnData)

	if err != nil {
		return 0, responseStr, err
	}

	return resp.StatusCode, returnData, nil
}

func StoreResult(result map[string]interface{}) (string, error) {
	dbName := os.Getenv("DBNAME")
	testData := os.Getenv("TESTDATA")
	client, err := database.GetDatabaseConnection()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	output, queryErr := client.Database(dbName).Collection(testData).InsertOne(ctx, result)
	if queryErr != nil {
		return "", err
	}
	objectId := output.InsertedID.(primitive.ObjectID)
	return objectId.String(), nil
}
