package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/database"
)

func getAll(c *gin.Context, repository database.Repository, dest interface{}) {
	if repository.FindAll(dest) {
		c.IndentedJSON(http.StatusOK, dest)
	} else {
		c.JSON(http.StatusInternalServerError, "Unable to process your request right now. Please check with system administrator.")
	}
}

func getByID(c *gin.Context, repository database.Repository, dest interface{}) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "IDs should be numerical values. Invalid ID received: "+c.Param("id"))
		return
	}

	if repository.FindById(dest, id) {
		c.IndentedJSON(http.StatusOK, dest)
	} else {
		c.JSON(http.StatusNotFound, fmt.Sprintf("{id: %d, message: \"Resource not found.\"}", id))
	}
}

func getByField(c *gin.Context, repository database.Repository, dest interface{}, query interface{}) {
	if repository.FindByField(dest, query) {
		c.IndentedJSON(http.StatusOK, dest)
	} else {
		c.JSON(http.StatusNotFound, fmt.Sprintf("{field: %s, message: \"Resource not found.\"}", query))
	}
}
