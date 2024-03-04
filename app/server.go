package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Serve() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", serveIndex)
	router.POST("/run", callModel)

	if err := router.Run(":31000"); err != nil {
		log.Fatal(err)
	}
}

func serveIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Message": "GoHTMX GPT",
	})
}
