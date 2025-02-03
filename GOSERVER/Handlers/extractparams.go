package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ExtractParams(c *gin.Context) map[string]string {
	defer PanicHandler(true, false, "Panic Enountered in ExtractParams function")

	params := make(map[string]string)

	// Extract query parameters for GET requests
	if c.Request.Method == "GET" {
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			} else {
				params[key] = ""
			}
		}
	}

	// Extract form parameters for POST requests
	if c.Request.Method == "POST" {
		c.Request.ParseMultipartForm(1024)
		for key, values := range c.Request.PostForm {
			if len(values) > 0 {
				params[key] = values[0]
			} else {
				params[key] = ""
			}
		}
	}

	if c.Request.Method == "PUT" {
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			} else {
				params[key] = ""
			}
		}

		if strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") {
			requestBody, readAllErr := io.ReadAll(c.Request.Body)    // storing request body in variable
			bodyReader := io.NopCloser(bytes.NewBuffer(requestBody)) // creating a reader with requestBody

			if readAllErr != nil {
				params["extract_error"] = readAllErr.Error()
			} else {
				postParams := make(map[string]interface{})
				decodeErr := json.NewDecoder(bodyReader).Decode(&postParams) // decoding the bodyReader
				if decodeErr != nil {
					params["extract_error"] = decodeErr.Error()
				} else {
					for key, values := range postParams {
						params[key] = InterfaceToString(values)
					}
				}
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // storing requestBody in c.Request.Body so that it can be decoded once again
		}
	}

	return params
}

func InterfaceToString(data interface{}) string {
	defer PanicHandler(true, false, "Panic Error encountered in InterfaceToString Function")
	result := ""
	switch x := data.(type) {
	case nil:
		result = ""
	case string:
		result = data.(string)
	case int:
		value := data.(int)
		result = strconv.Itoa(value)
	case int64:
		num := data.(int64)
		result = strconv.FormatInt(num, 10)
	case time.Time:
		date := data.(time.Time)
		result = date.Format("02-01-2006 15:04:05 -0700")
	case float64:
		value := data.(float64)
		result = strconv.FormatFloat(value, 'f', 0, 64)
	case map[string]interface{}:
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error:", err)
			return result
		}
		result = string(jsonData)
	case []byte:
		result = string(x)
	default:
		fmt.Printf("Value : %v", data)
		fmt.Println(data, " found as ", x)
	}
	return result
}

func PanicHandler(gchatAlertReq bool, mailRequired bool, panicMessage string) {
	if panicCheck := recover(); panicCheck != nil {
		panicErr := panicCheck.(error)
		additionalInfo := ""
		if len(panicMessage) > 0 {
			additionalInfo = panicMessage + "<br>The error encountered is :- " + panicErr.Error()
		}
		x := string(debug.Stack())
		message := ""
		if mailRequired {

			x = strings.ReplaceAll(x, "\n", "<br>")

			message = "<b style='font-size:1.2vw'>Panic error : </b><br><br><b>Stack Trace :- </b><br>" + x + "<br>" + "<b style='font-size:1.2vw'>" + additionalInfo + "</b>"

		}

		if gchatAlertReq {
			x = strings.ReplaceAll(x, "<br>", "\n")
			additionalInfo = panicMessage + "\n*The error is :- " + panicErr.Error() + "*"

			message = "*PANIC ERROR ENCOUNTERED !!!!*\n\n*Stack Trace :-*\n" + x + "\n" + additionalInfo + "\n*=========================================================================*\n\n"

		}
		fmt.Println(message)
		return
	}
}
