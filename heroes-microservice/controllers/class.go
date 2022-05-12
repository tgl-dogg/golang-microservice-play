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

// ClassHandler implements dependency injection for Repository. This controller needs no visibility to database connections.
type ClassHandler struct {
	repository database.Repository
}

// NewClassHandler constructs a new handler so we don't need to expose its internal fields.
func NewClassHandler(r database.Repository) ClassHandler {
	return ClassHandler{r}
}

// GetAll instances of this entity.
func (h *ClassHandler) GetAll(c *gin.Context) {
	getAll(c, h.repository, &[]heroes.Class{})
}

// GetByID the entity with the provided value in path parameter.
func (h *ClassHandler) GetByID(c *gin.Context) {
	getByID(c, h.repository, &heroes.Class{})
}

// GetByRole retrieve all entities whose role matches the provided value in path parameter.
func (h *ClassHandler) GetByRole(c *gin.Context) {
	role := heroes.Role(strings.ToLower(c.Param("role")))
	getByField(c, h.repository, &[]heroes.Class{}, &heroes.Class{Role: role})
}

// GetByProficiencies retrives all entities whose proficiencies match the parameters provided.
func (h *ClassHandler) GetByProficiencies(c *gin.Context) {
	var classes []heroes.Class
	proficiencies, queryParamNotEmpty := c.Request.URL.Query()["proficiencies"]

	if queryParamNotEmpty {
		// rawQuery := "SELECT * from classes c INNER JOIN class_proficiencies cp ON (cp.class_id = c.id) INNER JOIN proficiencies p ON (cp.proficiency_id = p.id) WHERE p.name IN ?"
		if err := h.repository.GetDB().Model(&classes).Distinct().Joins("INNER JOIN class_proficiencies cp ON (cp.class_id = id)").Joins("INNER JOIN proficiencies p ON (cp.proficiency_id = p.id)").Where("p.name IN ?", proficiencies).Find(&classes).Error; err != nil {
			log.Println("Error while executing getClassesByProficiencies: ", err)
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("{proficiencies: %s, message: \"Unable to process your request right now. Please check with system administrator.\"}", proficiencies))
			return
		}
	}

	c.IndentedJSON(http.StatusOK, classes)
}
