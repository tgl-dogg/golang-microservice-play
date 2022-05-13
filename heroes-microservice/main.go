package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tgl-dogg/golang-microservice-play/heroes-data"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/controllers"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/database"
)

func main() {
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

	skill := controllers.NewSkillHandler(repository)
	router.GET("/skills", skill.GetAll)
	router.GET("/skills/:id", skill.GetByID)
	router.GET("/skills/by-type/:type", skill.GetByType)
	router.GET("/skills/by-source/:source", skill.GetBySource)
}
