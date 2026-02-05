package repository

import (
	"github.com/zibianqu/eino_study/internal/model"
	"gorm.io/gorm"
)

type EntityRepository interface {
	Create(entity *model.Entity) error
	BatchCreate(entities []*model.Entity) error
	GetByDocID(docID string) ([]*model.Entity, error)
	GetByType(entityType string, offset, limit int) ([]*model.Entity, error)
	SearchByName(name string, offset, limit int) ([]*model.Entity, error)
	DeleteByDocID(docID string) error
}

type entityRepository struct {
	db *gorm.DB
}

func NewEntityRepository(db *gorm.DB) EntityRepository {
	return &entityRepository{db: db}
}

func (r *entityRepository) Create(entity *model.Entity) error {
	return r.db.Create(entity).Error
}

func (r *entityRepository) BatchCreate(entities []*model.Entity) error {
	return r.db.CreateInBatches(entities, 100).Error
}

func (r *entityRepository) GetByDocID(docID string) ([]*model.Entity, error) {
	var entities []*model.Entity
	err := r.db.Where("doc_id = ?", docID).Find(&entities).Error
	return entities, err
}

func (r *entityRepository) GetByType(entityType string, offset, limit int) ([]*model.Entity, error) {
	var entities []*model.Entity
	err := r.db.Where("entity_type = ?", entityType).
		Offset(offset).Limit(limit).
		Find(&entities).Error
	return entities, err
}

func (r *entityRepository) SearchByName(name string, offset, limit int) ([]*model.Entity, error) {
	var entities []*model.Entity
	err := r.db.Where("entity_name ILIKE ?", "%"+name+"%").
		Offset(offset).Limit(limit).
		Find(&entities).Error
	return entities, err
}

func (r *entityRepository) DeleteByDocID(docID string) error {
	return r.db.Where("doc_id = ?", docID).Delete(&model.Entity{}).Error
}