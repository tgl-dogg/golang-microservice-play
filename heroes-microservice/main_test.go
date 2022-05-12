package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	hero "github.com/tgl-dogg/golang-microservice-play/heroes-data"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/controllers"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var emptyRows = sqlmock.NewRows([]string{"id"})
var mockError = errors.New("Mock error.")

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

func Test_GetRaces_OK(t *testing.T) {
	r := gin.New()
	r.GET("/", getRaces)
	resp := emulateRequest(r, "/", http.StatusOK)

	var races []hero.Race
	decodeJSON(resp.Body, &races)

	if len(races) < 1 {
		t.Error("Invalid records found:", races)
	}
}

func Test_GetRaceByID_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	repository.GetDB()

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Human")
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	r := gin.New()
	r.GET("/:id", getRaceByID)
	resp := emulateRequest(r, "/1", http.StatusOK)

	var race hero.Race
	decodeJSON(resp.Body, &race)

	if race.Name != "Human" {
		t.Error("Invalid record found:", race)
	}

	shutdown(mock)
}

func Test_GetRaceByID_NOTFOUND(t *testing.T) {
	r := gin.New()
	r.GET("/:id", getRaceByID)
	emulateRequest(r, "/1000", http.StatusNotFound)
}

func Test_GetRaceByID_INVALID(t *testing.T) {
	invalidID := "98a11010-d019-11ec-9d64-0242ac120002"

	r := gin.New()
	r.GET("/:id", getRaceByID)
	resp := emulateRequest(r, "/"+invalidID, http.StatusBadRequest)

	body := resp.Body.String()
	if !strings.Contains(body, invalidID) {
		t.Error("Invalid response error:", body)
	}
}

func Test_GetRaceByRecommendedClasses_OK(t *testing.T) {
	r := gin.New()
	r.GET("/mock", getRacesByRecommendedClasses)
	resp := emulateRequest(r, "/mock?classes=thief&classes=warrior", http.StatusOK)

	var races []hero.Race
	decodeJSON(resp.Body, &races)

	if !isRacePresent(races, "Human") {
		t.Error("Expected Human race to be found.")
	} else if !isRacePresent(races, "Elf") {
		t.Error("Expected Elf race to be found.")
	} else if !isRacePresent(races, "Dwarf") {
		t.Error("Expected Dwarf race to be found.")
	}
}

func Test_GetRaceByRecommendedClasses_EMPTY(t *testing.T) {
	r := gin.New()
	r.GET("/mock", getRacesByRecommendedClasses)
	resp := emulateRequest(r, "/mock?classes=biruleibes", http.StatusOK)

	var races []hero.Race
	decodeJSON(resp.Body, &races)

	if len(races) > 0 {
		t.Error("Expected to found empty races set:", races)
	}
}

func isRacePresent(races []hero.Race, name string) bool {
	for i := range races {
		if races[i].Name == name {
			return true
		}
	}

	return false
}

func Test_GetClasses_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassController(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Warrior").AddRow(2, "Thief").AddRow(3, "Wizard")
	mock.ExpectQuery("SELECT (.+) FROM \"classes\"").WillReturnRows(rows)

	r := gin.New()
	r.GET("/", ch.GetAll)
	resp := emulateRequest(r, "/", http.StatusOK)

	var classes []hero.Class
	decodeJSON(resp.Body, &classes)

	if len(classes) != 3 {
		t.Error("Invalid records found:", classes)
	}

	shutdown(mock)
}

func Test_GetClasses_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassController(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"classes\"").WillReturnError(mockError)

	r := gin.New()
	r.GET("/", ch.GetAll)
	emulateRequest(r, "/", http.StatusInternalServerError)

	shutdown(mock)
}

func Test_GetClassByID_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassController(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Warrior")
	mock.ExpectQuery("SELECT (.+) FROM \"classes\" WHERE \"classes\".\"id\" = ? (.+)").WithArgs(1).WillReturnRows(rows)

	mock.ExpectQuery("SELECT (.+) FROM \"class_available_skills\" (.+)").WillReturnRows(mock.NewRows([]string{"id"}))
	mock.ExpectQuery("SELECT (.+) FROM \"class_proficiencies\" (.+)").WillReturnRows(mock.NewRows([]string{"id"}))
	mock.ExpectQuery("SELECT (.+) FROM \"class_starting_skills\" (.+)").WillReturnRows(mock.NewRows([]string{"id"}))

	r := gin.New()
	r.GET("/:id", ch.GetByID)
	resp := emulateRequest(r, "/1", http.StatusOK)

	var class hero.Class
	decodeJSON(resp.Body, &class)

	if class.Name != "Warrior" {
		t.Error("Invalid record found:", class)
	}

	shutdown(mock)
}

func Test_GetClassByID_INVALID(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassController(repository)

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
	ch := controllers.NewClassController(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"classes\" (.+)").WillReturnRows(emptyRows)

	r := gin.New()
	r.GET("/:id", ch.GetByID)
	emulateRequest(r, "/1", http.StatusNotFound)

	shutdown(mock)
}

func Test_GetClassByRole_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassController(repository)

	rows := mock.NewRows([]string{"id", "name", "role"}).AddRow(1, "Warrior", "fighter")
	mock.ExpectQuery("SELECT (.+) FROM \"classes\" WHERE \"classes\".\"role\" = ? (.+)").WithArgs("fighter").WillReturnRows(rows)

	r := gin.New()
	r.GET("/:role", ch.GetByRole)
	resp := emulateRequest(r, "/fighter", http.StatusOK)

	var classes []hero.Class
	decodeJSON(resp.Body, &classes)

	if classes[0].Role != "fighter" {
		t.Error("Invalid record found:", classes[0])
	}

	shutdown(mock)
}

func Test_GetClassByRole_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassController(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"classes\" (.+)").WillReturnError(mockError)

	r := gin.New()
	r.GET("/:role", ch.GetByRole)
	emulateRequest(r, "/malandro", http.StatusNotFound)

	shutdown(mock)
}

func Test_GetClassByProficiencies_OK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassController(repository)

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Warrior").AddRow(3, "Wizard")
	mock.ExpectQuery("SELECT (.+) FROM \"classes\"").WillReturnRows(rows)

	r := gin.New()
	r.GET("/mock", ch.GetByProficiencies)
	resp := emulateRequest(r, "/mock?proficiencies=complex_weapons&proficiencies=cast_magic", http.StatusOK)

	var classes []hero.Class
	decodeJSON(resp.Body, &classes)

	if len(classes) != 2 {
		t.Error("Invalid record found:", classes)
	}

	shutdown(mock)
}

func Test_GetClassByProficiencies_NOK(t *testing.T) {
	db, mock, repository := setup()
	defer db.Close()
	ch := controllers.NewClassController(repository)

	mock.ExpectQuery("SELECT (.+) FROM \"classes\"").WillReturnError(mockError)

	r := gin.New()
	r.GET("/mock", ch.GetByProficiencies)
	emulateRequest(r, "/mock?proficiencies=foresight", http.StatusNotFound)

	shutdown(mock)
}

func emulateRequest(r *gin.Engine, url string, expectedHttpStatus int) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != expectedHttpStatus {
		panic(fmt.Sprintf("HTTP request status code error. Expected: %d, found: %d", expectedHttpStatus, w.Code))
	}

	return w
}

func decodeJSON(r io.Reader, v any) {
	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		panic(err)
	}
}

// Sample races for mocking. Will be useful for unit testing later.
var mockRaces = []hero.Race{
	{
		ID:          1,
		Name:        "Human",
		Description: "We all understand the concept of a human. Lives in cities or whatever. Obs: plays with any classes.",
		BaseAttributes: hero.Attribute{
			Strength:     3,
			Agility:      3,
			Intelligence: 3,
			Willpower:    3,
		},
		StartingSkills:     []hero.Skill{},
		AvailableSkills:    []hero.Skill{},
		RecommendedClasses: mockClasses,
	},
	{
		ID:          2,
		Name:        "Elf",
		Description: "Pointy ears and snob noses. Cunning. Lives in forests.",
		BaseAttributes: hero.Attribute{
			Strength:     2,
			Agility:      4,
			Intelligence: 3,
			Willpower:    3,
		},
		StartingSkills:     []hero.Skill{},
		AvailableSkills:    []hero.Skill{},
		RecommendedClasses: []hero.Class{mockClasses[1]},
	},
	{
		ID:          3,
		Name:        "Dwarf",
		Description: "Small, strong and bearded. Likes mountains and steel.",
		BaseAttributes: hero.Attribute{
			Strength:     4,
			Agility:      2,
			Intelligence: 3,
			Willpower:    3,
		},
		StartingSkills:     []hero.Skill{mockSkills[0]},
		AvailableSkills:    []hero.Skill{},
		RecommendedClasses: []hero.Class{mockClasses[0]},
	},
}

// Sample classes for mocking. Will be useful for unit testing later.
var mockClasses = []hero.Class{
	{
		ID:          1,
		Name:        "Warrior",
		Description: "Powerful fighters that excel in tatical combat.",
		BonusAttributes: hero.Attribute{
			Strength:     1,
			Agility:      1,
			Intelligence: 0,
			Willpower:    0,
		},
		Role: hero.Fighter,
		Proficiencies: []hero.Proficiency{
			{ID: 1, Name: hero.SimpleWeapons},
			{ID: 2, Name: hero.ComplexWeapons},
		},
		StartingSkills:  []hero.Skill{},
		AvailableSkills: []hero.Skill{mockSkills[1]},
	},
	{
		ID:          2,
		Name:        "Thief",
		Description: "Elusive adventures capable of stealing things and pick locks without being noticed.",
		BonusAttributes: hero.Attribute{
			Strength:     0,
			Agility:      1,
			Intelligence: 1,
			Willpower:    0,
		},
		Role: hero.Dexterous,
		Proficiencies: []hero.Proficiency{
			{ID: 1, Name: hero.SimpleWeapons},
			{ID: 5, Name: hero.Pickpocket},
		},
		StartingSkills:  []hero.Skill{},
		AvailableSkills: []hero.Skill{},
	},
	{
		ID:          3,
		Name:        "Wizard",
		Description: "Arcane conjurers that can alter the tide of events with magic.",
		BonusAttributes: hero.Attribute{
			Strength:     0,
			Agility:      0,
			Intelligence: 1,
			Willpower:    1,
		},
		Role: hero.Spellcaster,
		Proficiencies: []hero.Proficiency{
			{ID: 1, Name: hero.SimpleWeapons},
			{ID: 3, Name: hero.CastMagic},
			{ID: 4, Name: hero.ReadMagic},
		},
		StartingSkills:  []hero.Skill{},
		AvailableSkills: []hero.Skill{mockSkills[2], mockSkills[3]},
	},
}

// Sample skills for mocking. Will be useful for unit testing later.
var mockSkills = []hero.Skill{
	{
		ID:                1,
		Name:              "Mountain Vigor",
		Description:       "You are immune to poisoning and can rol 3d6 when testing STR to resist fatigue.",
		Bonus:             "",
		Mana:              "",
		DifficultyType:    hero.Auto,
		Difficulty:        "",
		Activation:        hero.Passive,
		Source:            hero.FromRace,
		Type:              hero.Characteristic,
		LevelRequirement:  hero.None,
		SkillRequirements: []hero.Skill{},
		Observations:      "",
	}, {
		ID:                2,
		Name:              "War Cry",
		Description:       "You unleashe a fervorous scream that motiates your allies. You and them receive +1 in every roll until the end of the turn.",
		Bonus:             "This bonus is not cummulative.",
		Mana:              "10",
		DifficultyType:    hero.Auto,
		Difficulty:        "",
		Activation:        hero.Action,
		Source:            hero.FromClass,
		Type:              hero.Ability,
		LevelRequirement:  hero.None,
		SkillRequirements: []hero.Skill{},
		Observations:      "",
	}, {
		ID:                3,
		Name:              "Hellfire",
		Description:       "You engulf a 4m ground area in flames. Everyone making contact will take 10 damage (fire) and another 10 damage (fire) per subsequent round they remain there. Lasts for 3 rounds.",
		Bonus:             "Must be cast with a staff.",
		Mana:              "30",
		DifficultyType:    hero.Fixed,
		Difficulty:        "12",
		Activation:        hero.Action,
		Source:            hero.FromClass,
		Type:              hero.Spell,
		LevelRequirement:  hero.None,
		SkillRequirements: []hero.Skill{},
		Observations:      "",
	},
	{
		ID:                4,
		Name:              "Hellfire II",
		Description:       "You engulf a 4m diameter ground area in flames and it starts raining fire. Everyone inside will take 20 damage (fire) and another 20 damage (fire) per subsequent round they remain there. Lasts for 3 rounds.",
		Bonus:             "Must be cast with a staff.",
		Mana:              "40",
		DifficultyType:    hero.Fixed,
		Difficulty:        "12",
		Activation:        hero.Action,
		Source:            hero.FromClass,
		Type:              hero.Spell,
		LevelRequirement:  hero.Advanced,
		SkillRequirements: []hero.Skill{},
		Observations:      "",
	},
	{
		ID:                5,
		Name:              "Apprentice of [class]",
		Description:       "Immediately choose a different class than yours when acquiring this skill. You gain all of its proficiencies and can acquire its skills as of your own.",
		Bonus:             "",
		Mana:              "",
		DifficultyType:    hero.Auto,
		Difficulty:        "",
		Activation:        hero.Passive,
		Source:            hero.Base,
		Type:              hero.Technique,
		LevelRequirement:  hero.None,
		SkillRequirements: []hero.Skill{},
		Observations:      "You class is still considered to be your main class for any in-game purposes.",
	},
}
