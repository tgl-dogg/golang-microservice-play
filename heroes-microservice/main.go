package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tgl-dogg/golang-microservice-play/heroes-data"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/controllers"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	/*
	   Separar conexão num arquivo repository.go que recebe um DB no setup;
	   Separar uma classe com métodos de race, outra class e outra skills.
	   Teste de API usa classe mockada
	   Teste das classes usa DB mockado
	*/

	loadEnvFiles()
	setupDatabase()

	router := gin.Default()
	setupRoutes(router)
	router.Run("localhost:8080")
}

func loadEnvFiles() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured while loading .env file. Err: %s", err)
	}
}

func setupDatabase() {
	dbConnection := database.DBConnection{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		DBName:   os.Getenv("DATABASE_NAME"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
	}
	dbConnection.Setup()

	if os.Getenv("RUN_MIGRATIONS") == "true" {
		database.GetDB().AutoMigrate([]heroes.Skill{})
		database.GetDB().AutoMigrate([]heroes.Class{})
		database.GetDB().AutoMigrate([]heroes.Race{})
	}
}

func setupRoutes(router *gin.Engine) {
	repository := database.NewRepository(database.GetDB())

	router.GET("/races", getRaces)
	router.GET("/races/:id", getRaceByID)
	router.GET("/races/by-recommended-classes", getRacesByRecommendedClasses)

	class := controllers.NewClassController(repository)
	router.GET("/classes", class.GetAll)
	router.GET("/classes/:id", class.GetByID)
	router.GET("/classes/by-role/:role", class.GetByRole)
	router.GET("/classes/by-proficiencies", class.GetByProficiencies)

	router.GET("/skills", getSkills)
	router.GET("/skills/:id", getSkillByID)
	router.GET("/skills/by-type/:type", getSkillsByType)
	router.GET("/skills/by-source/:source", getSkillsBySource)
}

func getRaces(c *gin.Context) {
	var races []heroes.Race
	if findAll(c, &races) {
		c.IndentedJSON(http.StatusOK, races)
	}
}

func getRaceByID(c *gin.Context) {
	var race heroes.Race
	if findById(c, &race) {
		c.IndentedJSON(http.StatusOK, race)
	}
}

func getRacesByRecommendedClasses(c *gin.Context) {
	var races []heroes.Race
	queryClasses, queryParamNotEmpty := c.Request.URL.Query()["classes"]

	if queryParamNotEmpty {
		// Lowercasing params because SQL's IN clause is not case insensitive.
		for i := range queryClasses {
			queryClasses[i] = strings.ToLower(queryClasses[i])
		}

		if err := database.GetDB().Model(&races).Distinct().Preload("RecommendedClasses").Joins("INNER JOIN race_recommended_classes rc ON (rc.race_id = id)").Joins("INNER JOIN classes c ON (rc.class_id = c.id)").Where("LOWER(c.name) IN (?)", queryClasses).Find(&races).Error; err != nil {
			log.Println("Error while executing getRacesByRecommendedClasses: ", err)
			c.JSON(http.StatusNotFound, fmt.Sprintf("{classes: %s, message: \"Resource not found.\"}", queryClasses))
			return
		}
	}

	c.IndentedJSON(http.StatusOK, races)
}

func getClasses(c *gin.Context) {
	var classes []heroes.Class
	if findAll(c, &classes) {
		c.IndentedJSON(http.StatusOK, classes)
	}
}

func getClassByID(c *gin.Context) {
	var class heroes.Class
	if findById(c, &class) {
		c.IndentedJSON(http.StatusOK, class)
	}
}

func getClassesByRole(c *gin.Context) {
	var classes []heroes.Class
	role := heroes.Role(strings.ToLower(c.Param("role")))

	if findByField(c, &classes, &heroes.Class{Role: role}, "role", string(role)) {
		c.IndentedJSON(http.StatusOK, classes)
	}
}

func getClassesByProficiencies(c *gin.Context) {
	var classes []heroes.Class
	proficiencies, queryParamNotEmpty := c.Request.URL.Query()["proficiencies"]

	if queryParamNotEmpty {
		// rawQuery := "SELECT * from classes c INNER JOIN class_proficiencies cp ON (cp.class_id = c.id) INNER JOIN proficiencies p ON (cp.proficiency_id = p.id) WHERE p.name IN ?"
		if err := database.GetDB().Model(&classes).Distinct().Joins("INNER JOIN class_proficiencies cp ON (cp.class_id = id)").Joins("INNER JOIN proficiencies p ON (cp.proficiency_id = p.id)").Where("p.name IN ?", proficiencies).Find(&classes).Error; err != nil {
			log.Println("Error while executing getClassesByProficiencies: ", err)
			c.JSON(http.StatusNotFound, fmt.Sprintf("{proficiencies: %s, message: \"Resource not found.\"}", proficiencies))
			return
		}
	}

	c.IndentedJSON(http.StatusOK, classes)
}

func getSkills(c *gin.Context) {
	var skills []heroes.Skill
	if findAll(c, &skills) {
		c.IndentedJSON(http.StatusOK, skills)
	}
}

func getSkillByID(c *gin.Context) {
	var skill heroes.Skill
	if findById(c, &skill) {
		c.IndentedJSON(http.StatusOK, skill)
	}
}

func getSkillsByType(c *gin.Context) {
	var skills []heroes.Skill
	skillType := heroes.SkillType(strings.ToLower(c.Param("type")))

	if findByField(c, &skills, &heroes.Skill{Type: skillType}, "type", string(skillType)) {
		c.IndentedJSON(http.StatusOK, skills)
	}
}

func getSkillsBySource(c *gin.Context) {
	var skills []heroes.Skill
	source := heroes.Source(strings.ToLower(c.Param("source")))

	if findByField(c, &skills, &heroes.Skill{Source: source}, "source", string(source)) {
		c.IndentedJSON(http.StatusOK, skills)
	}
}

//*
func findAll(c *gin.Context, dest interface{}) bool {
	if err := database.GetDB().Find(dest).Error; err != nil {
		log.Println("Error while executing getAll: ", err)
		c.JSON(http.StatusInternalServerError, "Unable to process your request right now. Please check with system administrator.")
		return false
	}

	return true
}

func findById(c *gin.Context, dest interface{}) bool {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "IDs should be numerical values. Invalid ID received: "+c.Param("id"))
		return false
	}

	if err := database.GetDB().Preload(clause.Associations).First(dest, id).Error; err != nil {
		log.Println("Error while executing getByID: ", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, fmt.Sprintf("{id: %d, message: \"Resource not found.\"}", id))
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to process your request right now. Please check with system administrator.")
		}
		return false
	}

	return true
}

func findByField(c *gin.Context, dest interface{}, query interface{}, field string, value string) bool {
	if err := database.GetDB().Find(dest, query).Error; err != nil {
		log.Println("Error while executing findByField: ", err)
		c.JSON(http.StatusNotFound, fmt.Sprintf("{%s: %s, message: \"Resource not found.\"}", field, value))
		return false
	}

	return true
}

//*/
