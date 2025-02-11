package routing

import (
	handlers "github.com/ghostcode-sys/m/v2/Handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/writeFile", handlers.WriteFile)
	r.POST("/testCases", handlers.TestCases)
	r.GET("/loadHtml", handlers.LoadHtml)
	r.GET("/getData", handlers.FetchData)
	return r
}
