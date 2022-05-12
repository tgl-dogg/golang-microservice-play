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

type RaceHandler struct {
	repository database.Repository
}

func NewRaceHandler(r database.Repository) RaceHandler {
	return RaceHandler{r}
}

func (h *RaceHandler) GetAll(c *gin.Context) {
	getAll(c, h.repository, &[]heroes.Race{})
}

func (h *RaceHandler) GetByID(c *gin.Context) {
	getByID(c, h.repository, &heroes.Race{})
}

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
