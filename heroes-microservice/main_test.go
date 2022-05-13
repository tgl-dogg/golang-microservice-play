package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/tgl-dogg/golang-microservice-play/heroes-data"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/controllers"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var emptyRows = sqlmock.NewRows([]string{"id"})
var errMock = errors.New("just a mock error")

func setup() (db *sql.DB, mock sqlmock.Sqlmock, repository database.Repository) {
	// Open sqlmock connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("Error '%s' was not expected when opening a stub database connection.", err))
	}

	// Inject mocked connection into gormDB.
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to open gorm db, got error: %v", err))
	}

	// Pass our mocked database connection to the repository wrapper.
	repository = database.NewRepository(gormDB)
	return
}

func shutdown(mock sqlmock.Sqlmock) {
	// Make sure that all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		panic(fmt.Sprintf("There were unfulfilled expectations: %s", err))
	}
}

func Test_LoadEnvFile_OK(t *testing.T) {
	loadEnvFiles("../test.env")

	if os.Getenv("TEST_SUCCESSFUL") != "true" {
		t.Error("Expected env variables to be available.")
	}
}

func Test_LoadEnvFile_NOK(t *testing.T) {
	// This code should panic because it won't find any .env files to load.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic.")
		}
	}()

	loadEnvFiles("nobody.env")
}

func Test_SetupDatabase_NOK(t *testing.T) {
	// This code should panic because database connection string will be empty.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic.")
		}
	}()

	setupDatabase()
}

func Test_RunMigrations_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()

	countZero := sqlmock.NewRows([]string{"count"}).AddRow(0)
	successfulExec := sqlmock.NewResult(0, 0)
	mock.ExpectQuery("SELECT count(.+)").WillReturnRows(countZero)
	mock.ExpectExec("CREATE TABLE \"skills\" (.+)").WillReturnResult(successfulExec)
	mock.ExpectExec("CREATE TABLE \"skill_requirements\" (.+)").WillReturnResult(successfulExec)

	mock.ExpectQuery("SELECT count(.+)").WillReturnRows(countZero)
	mock.ExpectExec("CREATE TABLE \"classes\" (.+)").WillReturnResult(successfulExec)

	mock.ExpectQuery("SELECT count(.+)").WillReturnRows(countZero)
	mock.ExpectExec("CREATE TABLE \"races\" (.+)").WillReturnResult(successfulExec)

	os.Setenv("RUN_MIGRATIONS", "true")
	runMigrations(repository)

	shutdown(mock)
}

func Test_RunMigrations_SKIPPED(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()

	os.Setenv("RUN_MIGRATIONS", "false")
	runMigrations(repository)

	shutdown(mock)

}

func Test_RoutesRegisted_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()

	r := gin.New()
	setupRoutes(r, repository)

	routesMap := map[string]bool{
		"/races":                        false,
		"/races/:id":                    false,
		"/races/by-recommended-classes": false,
		"/classes":                      false,
		"/classes/:id":                  false,
		"/classes/by-role/:role":        false,
		"/classes/by-proficiencies":     false,
		"/skills":                       false,
		"/skills/:id":                   false,
		"/skills/by-type/:type":         false,
		"/skills/by-source/:source":     false,
	}

	for _, v := range r.Routes() {
		if _, ok := routesMap[v.Path]; ok {
			routesMap[v.Path] = true
		} else {
			t.Errorf("Unexpected route registered: %s", v.Path)
		}
	}

	for k, v := range routesMap {
		if !v {
			t.Errorf("Route expected but not registered: %s", k)
		}
	}

	shutdown(mock)
}

func Test_GetRaces_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewRaceHandler(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Human").AddRow(2, "Elf").AddRow(3, "Dwarf")
	mock.ExpectQuery("SELECT (.+) FROM \"races\"").WillReturnRows(rows)

	r := gin.New()
	r.GET("/", h.GetAll)
	resp := emulateRequest(r, "/", http.StatusOK)

	var races []heroes.Race
	decodeJSON(resp.Body, &races)

	if len(races) != 3 {
		t.Error("Invalid records found:", races)
	}

	shutdown(mock)
}

func Test_GetRaces_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewRaceHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"races\"").WillReturnError(errMock)

	r := gin.New()
	r.GET("/", h.GetAll)
	emulateRequest(r, "/", http.StatusInternalServerError)

	shutdown(mock)
}

func Test_GetRaceByID_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewRaceHandler(repository)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Human")
	mock.ExpectQuery("SELECT (.+) FROM \"races\" WHERE \"races\".\"id\" = ? (.+)").WithArgs(1).WillReturnRows(rows)

	mock.ExpectQuery("SELECT (.+) FROM \"race_available_skills\" (.+)").WillReturnRows(emptyRows)
	mock.ExpectQuery("SELECT (.+) FROM \"race_recommended_classes\" (.+)").WillReturnRows(emptyRows)
	mock.ExpectQuery("SELECT (.+) FROM \"race_starting_skills\" (.+)").WillReturnRows(emptyRows)

	r := gin.New()
	r.GET("/:id", h.GetByID)
	resp := emulateRequest(r, "/1", http.StatusOK)

	var race heroes.Race
	decodeJSON(resp.Body, &race)

	if race.Name != "Human" {
		t.Error("Invalid record found:", race)
	}

	shutdown(mock)
}

func Test_GetRaceByID_INVALID(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewRaceHandler(repository)

	invalidID := "98a11010-d019-11ec-9d64-0242ac120002"

	r := gin.New()
	r.GET("/:id", h.GetByID)
	resp := emulateRequest(r, "/"+invalidID, http.StatusBadRequest)

	body := resp.Body.String()
	if !strings.Contains(body, invalidID) {
		t.Error("Invalid response error:", body)
	}

	shutdown(mock)
}

func Test_GetRaceByID_NOTFOUND(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewRaceHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"races\" WHERE \"races\".\"id\" = ? (.+)").WillReturnRows(emptyRows)

	r := gin.New()
	r.GET("/:id", h.GetByID)
	emulateRequest(r, "/1000", http.StatusNotFound)

	shutdown(mock)
}

func Test_GetRaceByRecommendedClasses_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewRaceHandler(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Human").AddRow(3, "Dwarf")
	mock.ExpectQuery("SELECT (.+) FROM \"races\" (.+)").WillReturnRows(rows)
	mock.ExpectQuery("SELECT (.+) FROM \"race_recommended_classes\" (.+)").WillReturnRows(emptyRows)

	r := gin.New()
	r.GET("/mock", h.GetByRecommendedClasses)
	resp := emulateRequest(r, "/mock?classes=wizard&classes=warrior", http.StatusOK)

	var races []heroes.Race
	decodeJSON(resp.Body, &races)

	if len(races) != 2 {
		t.Error("Invalid records found:", races)
	}

	shutdown(mock)
}

func Test_GetRaceByRecommendedClasses_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewRaceHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"races\" (.+)").WillReturnError(errMock)

	r := gin.New()
	r.GET("/mock", h.GetByRecommendedClasses)
	emulateRequest(r, "/mock?classes=trickster", http.StatusInternalServerError)

	shutdown(mock)
}

func Test_GetClasses_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewClassHandler(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Warrior").AddRow(2, "Thief").AddRow(3, "Wizard")
	mock.ExpectQuery("SELECT (.+) FROM \"classes\"").WillReturnRows(rows)

	r := gin.New()
	r.GET("/", h.GetAll)
	resp := emulateRequest(r, "/", http.StatusOK)

	var classes []heroes.Class
	decodeJSON(resp.Body, &classes)

	if len(classes) != 3 {
		t.Error("Invalid records found:", classes)
	}

	shutdown(mock)
}

func Test_GetClasses_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewClassHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"classes\"").WillReturnError(errMock)

	r := gin.New()
	r.GET("/", h.GetAll)
	emulateRequest(r, "/", http.StatusInternalServerError)

	shutdown(mock)
}

func Test_GetClassByID_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassHandler(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Warrior")
	mock.ExpectQuery("SELECT (.+) FROM \"classes\" WHERE \"classes\".\"id\" = ? (.+)").WithArgs(1).WillReturnRows(rows)

	mock.ExpectQuery("SELECT (.+) FROM \"class_available_skills\" (.+)").WillReturnRows(emptyRows)
	mock.ExpectQuery("SELECT (.+) FROM \"class_proficiencies\" (.+)").WillReturnRows(emptyRows)
	mock.ExpectQuery("SELECT (.+) FROM \"class_starting_skills\" (.+)").WillReturnRows(emptyRows)

	r := gin.New()
	r.GET("/:id", ch.GetByID)
	resp := emulateRequest(r, "/1", http.StatusOK)

	var class heroes.Class
	decodeJSON(resp.Body, &class)

	if class.Name != "Warrior" {
		t.Error("Invalid record found:", class)
	}

	shutdown(mock)
}

func Test_GetClassByID_INVALID(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassHandler(repository)

	invalidID := "98a11010-d019-11ec-9d64-0242ac120002"
	r := gin.New()
	r.GET("/:id", ch.GetByID)
	resp := emulateRequest(r, "/"+invalidID, http.StatusBadRequest)

	body := resp.Body.String()
	if !strings.Contains(body, invalidID) {
		t.Error("Invalid response error:", body)
	}

	shutdown(mock)
}

func Test_GetClassByID_NOTFOUND(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"classes\" (.+)").WillReturnRows(emptyRows)

	r := gin.New()
	r.GET("/:id", ch.GetByID)
	emulateRequest(r, "/1", http.StatusNotFound)

	shutdown(mock)
}

func Test_GetClassByRole_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassHandler(repository)

	rows := mock.NewRows([]string{"id", "name", "role"}).AddRow(1, "Warrior", "fighter")
	mock.ExpectQuery("SELECT (.+) FROM \"classes\" WHERE \"classes\".\"role\" = ? (.+)").WithArgs("fighter").WillReturnRows(rows)

	r := gin.New()
	r.GET("/:role", ch.GetByRole)
	resp := emulateRequest(r, "/fighter", http.StatusOK)

	var classes []heroes.Class
	decodeJSON(resp.Body, &classes)

	if classes[0].Role != "fighter" {
		t.Error("Invalid record found:", classes[0])
	}

	shutdown(mock)
}

func Test_GetClassByRole_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"classes\" (.+)").WillReturnError(errMock)

	r := gin.New()
	r.GET("/:role", ch.GetByRole)
	emulateRequest(r, "/malandro", http.StatusInternalServerError)

	shutdown(mock)
}

func Test_GetClassByProficiencies_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassHandler(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Warrior").AddRow(3, "Wizard")
	mock.ExpectQuery("SELECT (.+) FROM \"classes\"").WillReturnRows(rows)

	r := gin.New()
	r.GET("/mock", ch.GetByProficiencies)
	resp := emulateRequest(r, "/mock?proficiencies=complex_weapons&proficiencies=cast_magic", http.StatusOK)

	var classes []heroes.Class
	decodeJSON(resp.Body, &classes)

	if len(classes) != 2 {
		t.Error("Invalid record found:", classes)
	}

	shutdown(mock)
}

func Test_GetClassByProficiencies_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"classes\"").WillReturnError(errMock)

	r := gin.New()
	r.GET("/mock", ch.GetByProficiencies)
	emulateRequest(r, "/mock?proficiencies=foresight", http.StatusInternalServerError)

	shutdown(mock)
}

func Test_GetSkills_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Mountain Vigor").AddRow(2, "War Cry").AddRow(3, "Hellfire")
	mock.ExpectQuery("SELECT (.+) FROM \"skills\"").WillReturnRows(rows)

	r := gin.New()
	r.GET("/", h.GetAll)
	resp := emulateRequest(r, "/", http.StatusOK)

	var skills []heroes.Skill
	decodeJSON(resp.Body, &skills)

	if len(skills) != 3 {
		t.Error("Invalid records found:", skills)
	}

	shutdown(mock)
}

func Test_GetSkills_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"skills\"").WillReturnError(errMock)

	r := gin.New()
	r.GET("/", h.GetAll)
	emulateRequest(r, "/", http.StatusInternalServerError)

	shutdown(mock)
}

func Test_GetSkillByID_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(2, "War Cry")
	mock.ExpectQuery("SELECT (.+) FROM \"skills\" WHERE \"skills\".\"id\" = ? (.+)").WithArgs(2).WillReturnRows(rows)

	mock.ExpectQuery("SELECT (.+) FROM \"skill_requirements\" (.+)").WillReturnRows(emptyRows)

	r := gin.New()
	r.GET("/:id", h.GetByID)
	resp := emulateRequest(r, "/2", http.StatusOK)

	var skill heroes.Skill
	decodeJSON(resp.Body, &skill)

	if skill.Name != "War Cry" {
		t.Error("Invalid record found:", skill)
	}

	shutdown(mock)
}

func Test_GetSkillByID_INVALID(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	invalidID := "98a11010-d019-11ec-9d64-0242ac120002"
	r := gin.New()
	r.GET("/:id", h.GetByID)
	resp := emulateRequest(r, "/"+invalidID, http.StatusBadRequest)

	body := resp.Body.String()
	if !strings.Contains(body, invalidID) {
		t.Error("Invalid response error:", body)
	}

	shutdown(mock)
}

func Test_GetSkillByID_NOTFOUND(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"skills\" (.+)").WillReturnRows(emptyRows)

	r := gin.New()
	r.GET("/:id", h.GetByID)
	emulateRequest(r, "/1", http.StatusNotFound)

	shutdown(mock)
}

func Test_GetSkillByType_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	rows := mock.NewRows([]string{"id", "name", "type"}).AddRow(3, "Hellfire", "spell").AddRow(4, "Hellfire II", "spell")
	mock.ExpectQuery("SELECT (.+) FROM \"skills\" WHERE \"skills\".\"type\" = ? (.+)").WithArgs("spell").WillReturnRows(rows)

	r := gin.New()
	r.GET("/:type", h.GetByType)
	resp := emulateRequest(r, "/spell", http.StatusOK)

	var skills []heroes.Skill
	decodeJSON(resp.Body, &skills)

	if skills[0].Type != "spell" {
		t.Error("Invalid record found:", skills[0])
	}

	shutdown(mock)
}

func Test_GetSkillByType_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"skills\"").WillReturnError(errMock)

	r := gin.New()
	r.GET("/:type", h.GetByType)
	emulateRequest(r, "/firula", http.StatusInternalServerError)

	shutdown(mock)
}

func Test_GetSkillBySource_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	rows := mock.NewRows([]string{"id", "name", "source"}).AddRow(1, "Mountain Vigor", "race")
	mock.ExpectQuery("SELECT (.+) FROM \"skills\" WHERE \"skills\".\"source\" = ? (.+)").WithArgs("race").WillReturnRows(rows)

	r := gin.New()
	r.GET("/:source", h.GetBySource)
	resp := emulateRequest(r, "/race", http.StatusOK)

	var skills []heroes.Skill
	decodeJSON(resp.Body, &skills)

	if skills[0].Source != "race" {
		t.Error("Invalid record found:", skills[0])
	}

	shutdown(mock)
}

func Test_GetSkillBySource_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	h := controllers.NewSkillHandler(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"skills\"").WillReturnError(errMock)

	r := gin.New()
	r.GET("/:source", h.GetBySource)
	emulateRequest(r, "/familiar", http.StatusInternalServerError)

	shutdown(mock)
}

func emulateRequest(r *gin.Engine, url string, expectedHTTPStatus int) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != expectedHTTPStatus {
		panic(fmt.Sprintf("HTTP request status code error. Expected: %d, found: %d", expectedHTTPStatus, w.Code))
	}

	return w
}

func decodeJSON(r io.Reader, v any) {
	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		panic(err)
	}
}
