package heroes

// Races represent the player's hero being, like Human or Elf.
type Race struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Strength           int     `json:"strength"`
	Agility            int     `json:"agility"`
	Intelligence       int     `json:"intelligence"`
	Willpower          int     `json:"willpower"`
	StartingSkills     []Skill `json:"starting_skills"`
	AvailableSkills    []Skill `json:"available_skills"`
	RecommendedClasses []Class `json:"recommendedClasses"`
}

// Classes represent how a hero is specialized, like Warrior or Wizard.
type Class struct {
	ID              int           `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Strength        int           `json:"strength"`
	Agility         int           `json:"agility"`
	Intelligence    int           `json:"intelligence"`
	Willpower       int           `json:"willpower"`
	Role            Role          `json:"role"`
	Proficiencies   []Proficiency `json:"proficiencies"`
	StartingSkills  []Skill       `json:"starting_skills"`
	AvailableSkills []Skill       `json:"available_skills"`
}

// Skills are heroes abilities. They can be either race or class skills, might require a determined level or previous skill knowledge and be techniques or spells.
type Skill struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Bonus            string    `json:"bonus"`
	Mana             string    `json:"mana"`
	Difficulty       string    `json:"difficulty"`
	Activation       string    `json:"activation"`
	Source           Source    `json:"source"`
	Type             SkillType `json:"type"`
	LevelRequirement int       `json:"level_requirement"`
	SkillRequirement []Skill   `json:"skill_requirement"`
}

// Proficiency represents natural abilities from classes, like being able to use complex weapons or cast magic. Classes come with two proficiencies, but might acquire more by multiclassing
type Proficiency string

const (
	SIMPLE_WEAPONS  Proficiency = "simple_weapons"
	COMPLEX_WEAPONS Proficiency = "complex_weapons"
	CAST_MAGIC      Proficiency = "cast_magic"
	READ_MAGIC      Proficiency = "read_magic"
	PICKPOKET       Proficiency = "pickpocket"
)

// Role represents overall class strategies: physical fighting, magical casting or dexterity usage. Classes usualy have only one role.
type Role string

const (
	FIGHTER     Role = "fighter"
	SPELLCASTER Role = "spellcaster"
	DEXTEROUS   Role = "dexterous"
)

// Skills can be learned as a racial feat, as a class ability, inherited from an acestor or be base (anyone can learn)
type Source string

const (
	BASE     Source = "base"
	RACE     Source = "race"
	CLASS    Source = "class"
	ANCESTOR Source = "ancestor"
)

// Skills have activation timing, which are actions (you perform in your turn), reactions (you performe after some condition happens) and passive (always active)
type Activation string

const (
	ACTION   Activation = "action"
	REACTION Activation = "reaction"
	PASSIVE  Activation = "passive"
)

// Skills can be simple abilities, characteristics (races usualy must be born with them), powerful tecniques (requires proficiency and/or some teaching) or spells (requires CAST_MAGIC proficiency, can be written in spellbooks)
type SkillType string

const (
	Ability        SkillType = "ability"
	Characteristic SkillType = "characteristic"
	Technique      SkillType = "technique"
	Spell          SkillType = "spell"
)
