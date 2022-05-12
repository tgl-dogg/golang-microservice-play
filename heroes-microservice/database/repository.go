package database

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db}
}

// GetDB gives direct access to gorm.DB capabilities.
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

func (r *Repository) FindAll(dest interface{}) bool {
	if err := r.db.Find(dest).Error; err != nil {
		log.Println("Error while executing getAll: ", err)
		return false
	}

	return true
}

func (r *Repository) FindById(dest interface{}, id uint64) bool {
	if err := r.db.Preload(clause.Associations).First(dest, id).Error; err != nil {
		log.Println("Error while executing getByID: ", err)

		return false
	}

	return true
}

func (r *Repository) FindByField(dest interface{}, query interface{}) bool {
	if err := r.db.Find(dest, query).Error; err != nil {
		log.Println("Error while executing findByField: ", err)
		return false
	}

	return true
}
