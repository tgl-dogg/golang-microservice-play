package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	hero "github.com/tgl-dogg/golang-microservice-play/heroes-data"
)

func init() {
	// Avoiding initialization loops in mocking.
	skills[3].SkillRequirement = append(skills[3].SkillRequirement, skills[2])
}

func main() {
	router := gin.Default()
	router.GET("/races", getRaces)
	router.GET("/races/:id", getRaceByID)
	router.GET("/races-by-recommended-classes", getRacesByRecommendedClasses)

	router.GET("/classes", getClasses)
	router.GET("/classes/:id", getClassByID)
	router.GET("/classes-by-role/:role", getClassesByRole)
	router.GET("/classes-by-proficiencies", getClassesByProficiencies)

	router.GET("/skills", getSkills)
	router.GET("/skills/:id", getSkillByID)
	router.GET("/skills-by-type/:type", getSkillsByType)
	router.GET("/skills-by-source/:source", getSkillsBySource)

	router.Run("localhost:8080")
}

func getRaces(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, races)
}

func getRaceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "IDs should be numerical values. Invalid ID received: "+c.Param("id"))
		return
	}

	for i := range races {
		if races[i].ID == id {
			c.IndentedJSON(http.StatusOK, races[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, fmt.Sprintf("{id: %d, message: \"Resource not found.\"}", id))
}

func getRacesByRecommendedClasses(c *gin.Context) {
	queryClasses, queryParamNotEmpty := c.Request.URL.Query()["classes"]
	results := make([]hero.Race, 0, len(races))

	if queryParamNotEmpty {
		// Not the best approach since it's O(n³), but will suffice for now before we add some real database.
		// Also, our RPG only has dozens of classes, no big deal.
		for i := range races {
		RECOMMENDED_CLASSES:
			for j := range races[i].RecommendedClasses {
				for k := range queryClasses {
					className := strings.ToLower(races[i].RecommendedClasses[j].Name)

					if strings.Contains(className, queryClasses[k]) {
						results = append(results, races[i])
						break RECOMMENDED_CLASSES // We only need to match a single class to consider the race.
					}
				}
			}
		}
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getClasses(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, classes)
}

func getClassByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "IDs should be numerical values. Invalid ID received: "+c.Param("id"))
		return
	}

	for i := range classes {
		if classes[i].ID == id {
			c.IndentedJSON(http.StatusOK, classes[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, fmt.Sprintf("{id: %d, message: \"Resource not found.\"}", id))
}

func getClassesByRole(c *gin.Context) {
	role := hero.Role(strings.ToLower(c.Param("role")))
	results := make([]hero.Class, 0, len(classes))

	for i := range classes {
		if classes[i].Role == role {
			results = append(results, classes[i])
		}
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getClassesByProficiencies(c *gin.Context) {
	proficiencies, queryParamNotEmpty := c.Request.URL.Query()["proficiencies"]
	results := make([]hero.Class, 0, len(classes))

	if queryParamNotEmpty {
		// Map helps deduplication, allows us some casting and provides a better way to write a "contains" feature.
		proficiencyMap := make(map[hero.Proficiency]string)
		for i := range proficiencies {
			proficiencyMap[hero.Proficiency(proficiencies[i])] = proficiencies[i]
		}

		for i := range classes {
			for j := range classes[i].Proficiencies {
				// Checks for "contains" proficiency in the proficiency map
				if _, ok := proficiencyMap[classes[i].Proficiencies[j]]; ok {
					results = append(results, classes[i])
					break // We only need to match a single proficiency to consider the class.
				}
			}
		}
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getSkills(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, skills)
}

func getSkillByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "IDs should be numerical values. Invalid ID received: "+c.Param("id"))
		return
	}

	for i := range skills {
		if skills[i].ID == id {
			c.IndentedJSON(http.StatusOK, skills[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, fmt.Sprintf("{id: %d, message: \"Resource not found.\"}", id))
}

func getSkillsByType(c *gin.Context) {
	skillType := hero.SkillType(strings.ToLower(c.Param("type")))
	results := make([]hero.Skill, 0, len(skills))

	for i := range skills {
		if skills[i].Type == skillType {
			results = append(results, skills[i])
		}
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getSkillsBySource(c *gin.Context) {
	source := hero.Source(strings.ToLower(c.Param("source")))
	results := make([]hero.Skill, 0, len(skills))

	for i := range skills {
		if skills[i].Source == source {
			results = append(results, skills[i])
		}
	}

	c.IndentedJSON(http.StatusOK, results)
}

// Sample races for mocking. Will be useful for unit testing later.
var races = []hero.Race{
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
		RecommendedClasses: classes,
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
		RecommendedClasses: []hero.Class{classes[1]},
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
		StartingSkills:     []hero.Skill{skills[0]},
		AvailableSkills:    []hero.Skill{},
		RecommendedClasses: []hero.Class{classes[0]},
	},
}

// Sample classes for mocking. Will be useful for unit testing later.
var classes = []hero.Class{
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
		Role:            hero.Fighter,
		Proficiencies:   []hero.Proficiency{hero.SimpleWeapons, hero.ComplexWeapons},
		StartingSkills:  []hero.Skill{},
		AvailableSkills: []hero.Skill{skills[1]},
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
		Role:            hero.Dexterous,
		Proficiencies:   []hero.Proficiency{hero.SimpleWeapons, hero.Pickpocket},
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
		Role:            hero.Spellcaster,
		Proficiencies:   []hero.Proficiency{hero.SimpleWeapons, hero.CastMagic, hero.ReadMagic},
		StartingSkills:  []hero.Skill{},
		AvailableSkills: []hero.Skill{skills[2], skills[3]},
	},
}

// Sample skills for mocking. Will be useful for unit testing later.
var skills = []hero.Skill{
	{
		ID:               1,
		Name:             "Mountain Vigor",
		Description:      "You are immune to poisoning and can rol 3d6 when testing STR to resist fatigue.",
		Bonus:            "",
		Mana:             "",
		DifficultyType:   hero.Auto,
		Difficulty:       "",
		Activation:       hero.Passive,
		Source:           hero.FromRace,
		Type:             hero.Characteristic,
		LevelRequirement: hero.None,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{},
	}, {
		ID:               2,
		Name:             "War Cry",
		Description:      "You unleashe a fervorous scream that motiates your allies. You and them receive +1 in every roll until the end of the turn.",
		Bonus:            "This bonus is not cummulative.",
		Mana:             "10",
		DifficultyType:   hero.Auto,
		Difficulty:       "",
		Activation:       hero.Action,
		Source:           hero.FromClass,
		Type:             hero.Ability,
		LevelRequirement: hero.None,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{},
	}, {
		ID:               3,
		Name:             "Hellfire",
		Description:      "You engulf a 4m ground area in flames. Everyone making contact will take 10 damage (fire) and another 10 damage (fire) per subsequent round they remain there. Lasts for 3 rounds.",
		Bonus:            "Must be cast with a staff.",
		Mana:             "30",
		DifficultyType:   hero.Fixed,
		Difficulty:       "12",
		Activation:       hero.Action,
		Source:           hero.FromClass,
		Type:             hero.Spell,
		LevelRequirement: hero.None,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{},
	},
	{
		ID:               4,
		Name:             "Hellfire II",
		Description:      "You engulf a 4m diameter ground area in flames and it starts raining fire. Everyone inside will take 20 damage (fire) and another 20 damage (fire) per subsequent round they remain there. Lasts for 3 rounds.",
		Bonus:            "Must be cast with a staff.",
		Mana:             "40",
		DifficultyType:   hero.Fixed,
		Difficulty:       "12",
		Activation:       hero.Action,
		Source:           hero.FromClass,
		Type:             hero.Spell,
		LevelRequirement: hero.Advanced,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{},
	},
	{
		ID:               5,
		Name:             "Apprentice of [class]",
		Description:      "Immediately choose a different class than yours when acquiring this skill. You gain all of its proficiencies and can acquire its skills as of your own.",
		Bonus:            "",
		Mana:             "",
		DifficultyType:   hero.Auto,
		Difficulty:       "",
		Activation:       hero.Passive,
		Source:           hero.Base,
		Type:             hero.Technique,
		LevelRequirement: hero.None,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{"You class is still considered to be your main class for any in-game purposes."},
	},
}
