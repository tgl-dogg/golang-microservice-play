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

// RaceHandler implements dependency injection for Repository. This controller needs no visibility to database connections.
type RaceHandler struct {
	repository database.Repository
}

// NewRaceHandler constructs a new handler so we don't need to expose its internal fields.
func NewRaceHandler(r database.Repository) RaceHandler {
	return RaceHandler{r}
}

// GetAll instances of this entity.
func (h *RaceHandler) GetAll(c *gin.Context) {
	getAll(c, h.repository, &[]heroes.Race{})
}

// GetByID the entity with the provided value in path parameter.
func (h *RaceHandler) GetByID(c *gin.Context) {
	getByID(c, h.repository, &heroes.Race{})
}

// GetByRecommendedClasses retrives all entities whose recommended classes match the parameters provided.
func (h *RaceHandler) GetByRecommendedClasses(c *gin.Context) {
	var races []heroes.Race
	queryClasses, queryParamNotEmpty := c.Request.URL.Query()["classes"]

	if queryParamNotEmpty {
		// Lowercasing params because SQL's IN clause is case sensitive.
		for i := range queryClasses {
			queryClasses[i] = strings.ToLower(queryClasses[i])
		}

		if err := h.repository.GetDB().Model(&races).Distinct().Preload("RecommendedClasses").Joins("INNER JOIN race_recommended_classes rc ON (rc.race_id = id)").Joins("INNER JOIN classes c ON (rc.class_id = c.id)").Where("LOWER(c.name) IN (?)", queryClasses).Find(&races).Error; err != nil {
			log.Println("Error while executing getRacesByRecommendedClasses: ", err)
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("{classes: %s, message: \"Unable to process your request right now. Please check with system administrator.\"}", queryClasses))
			return
		}
	}

	c.IndentedJSON(http.StatusOK, races)
}
