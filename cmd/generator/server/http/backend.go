package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

var r = gin.Default()
var state = map[string]interface{}{}



func MainHandler() {
	TestData[0].MakeRequest()
	r.GET("/todos", func(c *gin.Context) {
		c.JSON(200, state)
	})
	r.GET("/todos/:id", func(c *gin.Context) {
		todo := state[c.Param("id")]
		c.JSON(200, todo)
	})
	r.POST("/todos", func(c *gin.Context) {
		var todo interface{}
		c.BindJSON(&todo)
		state[strconv.Itoa(len(state))] = todo
		c.JSON(200, todo)
	})
	r.PUT("/todos/:id", func(c *gin.Context) {
		todo := state[c.Param("id")]
		c.JSON(200, todo)
	})
	r.PATCH("/todos/:id", func(c *gin.Context) {
		todo := state[c.Param("id")]
		c.JSON(200, todo)
	})
	r.DELETE("/todos/:id", func(c *gin.Context) {
		todo := state[c.Param("id")]
		delete(state, c.Param("id"))
		c.JSON(200, todo)
	})
	//	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
