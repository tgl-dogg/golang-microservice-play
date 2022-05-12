package database

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository implements dependency injection for database connection.
type Repository struct {
	db *gorm.DB
}

// NewRepository constructs a new Repository so we don't need to expose Repository's internal fields.
func NewRepository(db *gorm.DB) Repository {
	return Repository{db}
}

// GetDB gives direct access to gorm.DB capabilities.
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

// FindAll is an abstraction of gorm.Find. Searches all records of the desired interface.
func (r *Repository) FindAll(dest interface{}) bool {
	if err := r.db.Find(dest).Error; err != nil {
		log.Println("Error while executing getAll: ", err)
		return false
	}

	return true
}

// FindByID is an abstraction of gorm.Find using primary key. Searches desired interface using provided primary key
func (r *Repository) FindByID(dest interface{}, id uint64) bool {
	if err := r.db.Preload(clause.Associations).First(dest, id).Error; err != nil {
		log.Println("Error while executing getByID: ", err)

		return false
	}

	return true
}

// FindByField is an abstraction of gorm.Find. Finds the desired interface applying the provided query parameter.
func (r *Repository) FindByField(dest interface{}, query interface{}) bool {
	if err := r.db.Find(dest, query).Error; err != nil {
		log.Println("Error while executing findByField: ", err)
		return false
	}

	return true
}
