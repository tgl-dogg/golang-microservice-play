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
	// SimpleWeapons allows usage of small cold weapons, such as daggers, shortswords, handaxes, bows and crossbows.
	SimpleWeapons Proficiency = "simple_weapons"
	// ComplexWeapons allows usage of bigger cold weapons, such as longswords, greataxes, lances, warbows and heavy crossbows.
	ComplexWeapons Proficiency = "complex_weapons"
	// CastMagic allows spellcasting, such as Wizard's Fireball or Hellfire.
	CastMagic Proficiency = "cast_magic"
	// ReadMagic allows reading magically engraved itens, such as spellbooks, runes or enchanted weapons.
	ReadMagic Proficiency = "read_magic"
	// Pickpocket allows picking locks, disarming traps and stealing from unsuspecting pockets.
	Pickpocket Proficiency = "pickpocket"
)

// Role represents overall class strategies: physical fighting, magical casting or dexterity usage. Classes usualy have only one role.
type Role string

const (
	// Fighter uses melee weapons and has greater overall endurance with high damage output. Excels at strength and defense, has high agility and hit points.
	Fighter Role = "fighter"
	// Spellcaster can be versatile due its access to spells. Can sustain damage, provide support or be a jack-of-all-trades. Excels at intelligence and has high willpower and mana.
	Spellcaster Role = "spellcaster"
	// Dexterous are cunning and deceiving, may fight at distance, work treacherously or avoid being noticed at all. Excels at agility and dodge, has high intelligence and balanced attributes.
	Dexterous Role = "dexterous"
)

// DifficultyType represents skills different difficulties, which are a target number that must be achieved with hero attribute + dice rolling to be performed.
type DifficultyType string

const (
	// Auto skills are always active (if passive) or automatically used upon activation, having no difficult target like Warrior's War Cry.
	Auto DifficultyType = "auto"
	// Fixed skills have a fixed target number as difficulty, like 12 for Wizard's Fireball or Hellfire.
	Fixed DifficultyType = "fixed"
	// Variable difficulty depends on player roleplaying choices, like trying to levitate a small rock or a cow.
	Variable DifficultyType = "variable"
	// TargetPlus difficulty is set upon a target value (like opponent's defense or dodge) with a modifier (can be positive, negative or zero).
	TargetPlus DifficultyType = "target_plus"
)

// Activation represents skills timing, like activating in your turn, in response to something or being always active.
type Activation string

const (
	// Action is performed during your turn, such as Warrior's War Cry or Wizard's Hellfire.
	Action Activation = "action"
	// Reaction skills activate after some precondition happens, like response to taking damage.
	Reaction Activation = "reaction"
	// Passive skills are always active, like Dwarf's Mountain Vigor.
	Passive Activation = "passive"
)

// Source is where a skills can be learned from.
type Source string

const (
	// Base skills can be leartn by anyone. Skill requirements must still be met.
	Base Source = "base"
	// FromRace can only be accessed by members of a determined race.
	FromRace Source = "race"
	// FromClass can only be accessed by members of a determined class.
	FromClass Source = "class"
	// FromAncestor must be learnt from your ancestral inheritance.
	FromAncestor Source = "ancestor"
)

// SkillType represents skills categories, like abilities, characteristics (races are usualy born with them), powerful tecniques (requires proficiency and/or some teaching)
// or spells (requires CAST_MAGIC proficiency, can be written in spellbooks).
type SkillType string

const (
	// Ability is a simple skill type, with nothing special about it.
	Ability SkillType = "ability"
	// Characteristic is usualy a passive or racial feat. Characteristics become your hero's way of being and might change attributes or physical appearance (like having four arms).
	Characteristic SkillType = "characteristic"
	// Technique usually requires proficiency and can be learnt by spending skill points or training with a mentor in-game.
	Technique SkillType = "technique"
	// Spell requires cast_spell proficiency and can be executed by memory (when you acquire the skill) or from a spellbook (if you have studied it beforehand).
	Spell SkillType = "spell"
)

// LevelRequirement is a skills dependency model, which can be advanced (level 5 or above), master (level 10 or above) or initial (must be level one to acquire).
type LevelRequirement string

const (
	// None means you can learn the skill at any level.
	None LevelRequirement = "none"
	// Advanced skills must be learnt at level 5 or above. These skills are powerful game changers.
	Advanced LevelRequirement = "advanced"
	// Master skills must be learnt at level 10 or above. Classes usualy have one or two master skills at most, as they are ultimate skills.
	Master LevelRequirement = "master"
	// Initial skills must be learnt at level 1, i.e. when you first create your hero sheet. This usualy includes ancestor skills and racial feats (like having four arms).
	Initial LevelRequirement = "initial"
)
