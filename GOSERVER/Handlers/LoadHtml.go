package handlers

import "github.com/gin-gonic/gin"

func LoadHtml(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}
