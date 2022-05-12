package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tgl-dogg/golang-microservice-play/heroes-data"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/database"
)

type ClassHandler struct {
	repository database.Repository
}

func NewClassHandler(r database.Repository) ClassHandler {
	return ClassHandler{r}
}

func (h *ClassHandler) GetAll(c *gin.Context) {
	getAll(c, h.repository, &[]heroes.Class{})
}

func (h *ClassHandler) GetByID(c *gin.Context) {
	getByID(c, h.repository, &heroes.Class{})
}

func (h *ClassHandler) GetByRole(c *gin.Context) {
	role := heroes.Role(strings.ToLower(c.Param("role")))
	getByField(c, h.repository, &[]heroes.Class{}, &heroes.Class{Role: role})
}

func (h *ClassHandler) GetByProficiencies(c *gin.Context) {
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
