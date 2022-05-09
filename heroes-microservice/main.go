package main

import (
	"strings"
	hero "tgl-dogg/heroes-data"

	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// Avoiding initalization loops in mocking.
	skills[3].SkillRequirement = append(skills[3].SkillRequirement, skills[2])
}

func main() {
	router := gin.Default()
	router.GET("/races", getRaces)
	//router.GET("/races-by-recommended-class", getRacesByRecommendedClass)

	router.GET("/classes", getClasses)
	router.GET("/classes-by-role/:role", getClassesByRole)
	//router.GET("/classes-by-proficiencies", getClassesByProficiencies)

	router.GET("/skills", getSkills)
	//router.GET("/skills-by-type", getSkillsByType)
	//router.GET("/skills-by-source", getSkillsBySource)

	router.Run("localhost:8080")
}

func getRaces(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, races)
}

func getClasses(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, classes)
}

func getClassesByRole(c *gin.Context) {
	role := strings.ToLower(c.Param("role"))
	classesByRole := make([]hero.Class, 0, len(classes))

	for i := range classes {
		if string(classes[i].Role) == role {
			classesByRole = append(classesByRole, classes[i])
		}
	}

	c.IndentedJSON(http.StatusOK, classesByRole)
}

func getSkills(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, skills)
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
		Role:            hero.FIGHTER,
		Proficiencies:   []hero.Proficiency{hero.SIMPLE_WEAPONS, hero.COMPLEX_WEAPONS},
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
		Role:            hero.DEXTEROUS,
		Proficiencies:   []hero.Proficiency{hero.SIMPLE_WEAPONS, hero.PICKPOKET},
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
		Role:            hero.SPELLCASTER,
		Proficiencies:   []hero.Proficiency{hero.SIMPLE_WEAPONS, hero.CAST_MAGIC},
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
		Difficulty:       "",
		Activation:       hero.PASSIVE,
		Source:           hero.RACE,
		Type:             hero.CHARACTERISTIC,
		LevelRequirement: hero.NONE,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{},
	}, {
		ID:               2,
		Name:             "War Cry",
		Description:      "You unleashe a fervorous scream that motiates your allies. You and them receive +1 in every roll until the end of the turn.",
		Bonus:            "This bonus is not cummulative.",
		Mana:             "10",
		Difficulty:       "",
		Activation:       hero.ACTION,
		Source:           hero.CLASS,
		Type:             hero.ABILITY,
		LevelRequirement: hero.NONE,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{},
	}, {
		ID:               3,
		Name:             "Hellfire",
		Description:      "You engulf a 4m ground area in flames. Everyone making contact will take 10 damage (fire) and another 10 damage (fire) per subsequent round they remain there. Lasts for 3 rounds.",
		Bonus:            "Must be cast with a staff.",
		Mana:             "30",
		Difficulty:       "12",
		Activation:       hero.ACTION,
		Source:           hero.CLASS,
		Type:             hero.SPELL,
		LevelRequirement: hero.NONE,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{},
	},
	{
		ID:               4,
		Name:             "Hellfire II",
		Description:      "You engulf a 4m diameter ground area in flames and create a rain of fire. Everyone inside will take 20 damage (fire) and another 20 damage (fire) per subsequent round they remain there. Lasts for 3 rounds.",
		Bonus:            "Must be cast with a staff.",
		Mana:             "40",
		Difficulty:       "12",
		Activation:       hero.ACTION,
		Source:           hero.CLASS,
		Type:             hero.SPELL,
		LevelRequirement: hero.ADVANCED,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{},
	},
	{
		ID:               5,
		Name:             "Apprentice of [class]",
		Description:      "Immediately choose a different class than yours when acquiring this skill. You gain all of its proficiencies and can acquire its skills as of your own.",
		Bonus:            "",
		Mana:             "",
		Difficulty:       "",
		Activation:       hero.PASSIVE,
		Source:           hero.BASE,
		Type:             hero.TECHNIQUE,
		LevelRequirement: hero.NONE,
		SkillRequirement: []hero.Skill{},
		Observations:     []string{"You are still considered your main class for any in-game purposes."},
	},
}
