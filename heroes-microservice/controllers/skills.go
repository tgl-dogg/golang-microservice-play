package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tgl-dogg/golang-microservice-play/heroes-data"
	"github.com/tgl-dogg/golang-microservice-play/heroes-microservice/database"
)

// SkillHandler implements dependency injection for Repository. This controller needs no visibility to database connections.
type SkillHandler struct {
	repository database.Repository
}

// NewSkillHandler constructs a new handler so we don't need to expose its internal fields.
func NewSkillHandler(r database.Repository) SkillHandler {
	return SkillHandler{r}
}

// GetAll instances of this entity.
func (h *SkillHandler) GetAll(c *gin.Context) {
	getAll(c, h.repository, &[]heroes.Skill{})
}

// GetByID the entity with the provided value in path parameter.
func (h *SkillHandler) GetByID(c *gin.Context) {
	getByID(c, h.repository, &heroes.Skill{})
}

// GetByType retrieve all entities whose source matches the provided value in path parameter.
func (h *SkillHandler) GetByType(c *gin.Context) {
	skillType := heroes.SkillType(strings.ToLower(c.Param("type")))
	getByField(c, h.repository, &[]heroes.Skill{}, &heroes.Skill{Type: skillType})
}

// GetBySource retrieve all entities whose source matches the provided value in path parameter.
func (h *SkillHandler) GetBySource(c *gin.Context) {
	source := heroes.Source(strings.ToLower(c.Param("source")))
	getByField(c, h.repository, &[]heroes.Skill{}, &heroes.Skill{Source: source})
}
