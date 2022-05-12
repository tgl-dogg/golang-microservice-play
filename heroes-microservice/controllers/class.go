package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tgl-dogg/golang-microservice-play/heroes-data"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/database"
)

type ClassController struct {
	repository database.Repository
}

func NewClassController(r database.Repository) ClassController {
	return ClassController{r}
}

func (h *ClassController) GetAll(c *gin.Context) {
	var classes []heroes.Class
	if h.repository.FindAll(&classes) {
		c.IndentedJSON(http.StatusOK, classes)
	} else {
		c.JSON(http.StatusInternalServerError, "Unable to process your request right now. Please check with system administrator.")
	}
}

func (h *ClassController) GetByID(c *gin.Context) {
	var class heroes.Class

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "IDs should be numerical values. Invalid ID received: "+c.Param("id"))
		return
	}

	if h.repository.FindById(&class, id) {
		c.IndentedJSON(http.StatusOK, class)
	} else {
		c.JSON(http.StatusNotFound, fmt.Sprintf("{id: %d, message: \"Resource not found.\"}", id))
	}
}

func (h *ClassController) GetByRole(c *gin.Context) {
	var classes []heroes.Class
	role := heroes.Role(strings.ToLower(c.Param("role")))

	// "role", string(role)
	if h.repository.FindByField(&classes, &heroes.Class{Role: role}) {
		c.IndentedJSON(http.StatusOK, classes)
	} else {
		c.JSON(http.StatusNotFound, fmt.Sprintf("{%s: %s, message: \"Resource not found.\"}", "role", role))
	}
}

func (h *ClassController) GetByProficiencies(c *gin.Context) {
	var classes []heroes.Class
	proficiencies, queryParamNotEmpty := c.Request.URL.Query()["proficiencies"]

	if queryParamNotEmpty {
		// rawQuery := "SELECT * from classes c INNER JOIN class_proficiencies cp ON (cp.class_id = c.id) INNER JOIN proficiencies p ON (cp.proficiency_id = p.id) WHERE p.name IN ?"
		if err := h.repository.GetDB().Model(&classes).Distinct().Joins("INNER JOIN class_proficiencies cp ON (cp.class_id = id)").Joins("INNER JOIN proficiencies p ON (cp.proficiency_id = p.id)").Where("p.name IN ?", proficiencies).Find(&classes).Error; err != nil {
			log.Println("Error while executing getClassesByProficiencies: ", err)
			c.JSON(http.StatusNotFound, fmt.Sprintf("{proficiencies: %s, message: \"Resource not found.\"}", proficiencies))
			return
		}
	}

	c.IndentedJSON(http.StatusOK, classes)
}
