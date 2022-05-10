package heroes

// Race represents the player's hero being, like Human or Elf. They have base attributes for the character and learn racial skills.
type Race struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	BaseAttributes     Attribute `json:"base_attributes"`
	StartingSkills     []Skill   `json:"starting_skills"`
	AvailableSkills    []Skill   `json:"available_skills"`
	RecommendedClasses []Class   `json:"recommendedClasses"`
}

// Class represents how a hero is specialized, like Warrior or Wizard. They give bonus attributes depending on their skillset.
type Class struct {
	ID              int           `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	BonusAttributes Attribute     `json:"bonus_attributes"`
	Role            Role          `json:"role"`
	Proficiencies   []Proficiency `json:"proficiencies"`
	StartingSkills  []Skill       `json:"starting_skills"`
	AvailableSkills []Skill       `json:"available_skills"`
}

// Skill is a hero ability.
// They can be either race or class skills, might be techniques or spells and require a minimum level or previous skill knowledge.
type Skill struct {
	ID               int              `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Bonus            string           `json:"bonus"`
	Mana             string           `json:"mana"`
	DifficultyType   DifficultyType   `json:"difficulty_type"`
	Difficulty       string           `json:"difficulty"`
	Activation       Activation       `json:"activation"`
	Source           Source           `json:"source"`
	Type             SkillType        `json:"type"`
	LevelRequirement LevelRequirement `json:"level_requirement"`
	SkillRequirement []Skill          `json:"skill_requirement"`
	Observations     []string         `json:"observations"`
}

// Attribute is a hero measurement of power. Heroes have strength (physical power), agility (velocity and dexterity),
// intelligence (smartness and cast magic) and overall willpower.
type Attribute struct {
	Strength     int `json:"strength"`
	Agility      int `json:"agility"`
	Intelligence int `json:"intelligence"`
	Willpower    int `json:"willpower"`
}

// Proficiency represents natural abilities from classes, like being able to use complex weapons or cast magic.
// Classes come with two proficiencies, but might acquire more by multiclassing.
type Proficiency string

const (
	// Case use small cold weapons, such as daggers, shortswords, handaxes, bows and crossbows.
	SimpleWeapons Proficiency = "simple_weapons"
	// Can use bigger cold weapons, such as longswords, greataxes, lances, warbows and heavy crossbows.
	ComplexWeapons Proficiency = "complex_weapons"
	// Can cast spells, such as Fireball or Hellfire.
	CastMagic Proficiency = "cast_magic"
	// Can read magically engraved itens, such as spellbooks, runes or enchanted weapons.
	ReadMagic Proficiency = "read_magic"
	// Can pick locks, disarm traps and steal from unsuspecting pockets.
	Pickpocket Proficiency = "pickpocket"
)

// Role represents overall class strategies: physical fighting, magical casting or dexterity usage. Classes usualy have only one role.
type Role string

const (
	// Fights with melee weapons and has greater overall endurance with high damage output. Excels at strength and defense, has high agility and hit points.
	Fighter Role = "fighter"
	// Versatile due its access to spells, can sustain damage, provide support or be a jack-of-all-trades. Excels at intelligence and has high willpower and mana.
	Spellcaster Role = "spellcaster"
	// Cunning and deceiving, may fight at distance, work treacherously or avoid being noticed at all. Excels at agility and dodge, has high intelligence and balanced attributes.
	Dexterous Role = "dexterous"
)

// DifficultyType represents skills different difficulties, which are a target number that must be achieved with hero attribute + dice rolling to be performed.
type DifficultyType string

const (
	// Skill is always active (if passive) or automatically used upon activation timing without requirent a test, like Warrior's War Cry.
	Auto DifficultyType = "auto"
	// Skill have a fixed target number to be performed, like 12 for Wizard's Fireball or Hellfire.
	Fixed DifficultyType = "fixed"
	// Difficulty depends on player roleplaying choice, like trying to levitate a small rock or a cow.
	Variable DifficultyType = "variable"
	// Difficult is set upon a target value (like opponent's defense or dodge) with a modifier (can be positive, negative or zero).
	TargetPlus DifficultyType = "target_plus"
)

// Activation represents skills timing, like activating in your turn, in response to something or being always active.
type Activation string

const (
	// You perform during your turn, such as Warrior's War Cry.
	Action Activation = "action"
	// Skill activate after some precondition happens, like response to taking damage.
	Reaction Activation = "reaction"
	// Always active, like Dwarf's Mountain Vigor.
	Passive Activation = "passive"
)

// Source is where a skills can be learned from.
type Source string

const (
	// Anyone can learn. Skill requirements must still be met.
	Base Source = "base"
	// Can only be accessed by members of a determined race.
	FromRace Source = "race"
	// Can only be accessed by members of a determined class.
	FromClass Source = "class"
	// Must be learnt from your ancestral inheritance.
	FromAncestor Source = "ancestor"
)

// SkillType represents skills categories, like abilities, characteristics (races are usualy born with them), powerful tecniques (requires proficiency and/or some teaching)
// or spells (requires CAST_MAGIC proficiency, can be written in spellbooks).
type SkillType string

const (
	// Simple skill type, has nothing special.
	Ability SkillType = "ability"
	// Usualy passives or racial feats. Characteristics become your hero's way of being and might change attributes or physical appearance (like having four arms).
	Characteristic SkillType = "characteristic"
	// Techniques usually require proficiency and can be learnt by spending skill points or training with a mentor in-game.
	Technique SkillType = "technique"
	// Spells require cast_spell proficiency and can be executed by memory (when you acquire the skill) or from a spellbook (if you have studied it beforehand).
	Spell SkillType = "spell"
)

// LevelRequirement is a skills dependency model, which can be advanced (level 5 or above), master (level 10 or above) or initial (must be level one to acquire).
type LevelRequirement string

const (
	// No level requirement. Can be learnt at any time, provided your have spare skill points and access to it.
	None LevelRequirement = "none"
	// Must be level 5 or above to learn. These skills are powerful game changers.
	Advanced LevelRequirement = "advanced"
	// Must be level 10 or above to learn. Classes usualy have one or two master skills at most, as they are ultimate skills.
	Master LevelRequirement = "master"
	// Must be learnt at level 1, i.e. when you first create your hero sheet.
	Initial LevelRequirement = "initial"
)
