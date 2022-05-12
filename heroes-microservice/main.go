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

	race := controllers.NewRaceHandler(repository)
	router.GET("/races", race.GetAll)
	router.GET("/races/:id", race.GetByID)
	router.GET("/races/by-recommended-classes", race.GetByRecommendedClasses)

	class := controllers.NewClassHandler(repository)
	router.GET("/classes", class.GetAll)
	router.GET("/classes/:id", class.GetByID)
	router.GET("/classes/by-role/:role", class.GetByRole)
	router.GET("/classes/by-proficiencies", class.GetByProficiencies)

	router.GET("/skills", getSkills)
	router.GET("/skills/:id", getSkillByID)
	router.GET("/skills/by-type/:type", getSkillsByType)
	router.GET("/skills/by-source/:source", getSkillsBySource)
}

func getSkills(c *gin.Context) {
	var skills []heroes.Skill
	if findAll(c, &skills) {
		c.IndentedJSON(http.StatusOK, skills)
	}
}

func getSkillByID(c *gin.Context) {
	var skill heroes.Skill
	if findByID(c, &skill) {
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

func findByID(c *gin.Context, dest interface{}) bool {
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
