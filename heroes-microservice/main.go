package main

import (
	hero "tgl-dogg/heroes-data"

	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/races", getRaces)

	router.Run("localhost:8080")
}

func getRaces(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, races)
}

// sample
var races = []hero.Race{
	{
		ID:                 1,
		Name:               "Human",
		Description:        "We all understand the concept of a human.",
		Strength:           3,
		Agility:            3,
		Intelligence:       3,
		Willpower:          3,
		StartingSkills:     []hero.Skill{},
		AvailableSkills:    []hero.Skill{},
		RecommendedClasses: []hero.Class{},
	},
	{
		ID:                 2,
		Name:               "Elf",
		Description:        "Pointy ears and snob noses. Live in forests.",
		Strength:           2,
		Agility:            4,
		Intelligence:       3,
		Willpower:          3,
		StartingSkills:     []hero.Skill{},
		AvailableSkills:    []hero.Skill{},
		RecommendedClasses: []hero.Class{},
	},
	{
		ID:                 3,
		Name:               "Dwarf",
		Description:        "Small and strong, like montains and steel.",
		Strength:           4,
		Agility:            2,
		Intelligence:       3,
		Willpower:          3,
		StartingSkills:     []hero.Skill{},
		AvailableSkills:    []hero.Skill{},
		RecommendedClasses: []hero.Class{},
	},
}
