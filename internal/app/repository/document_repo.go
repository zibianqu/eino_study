package repository

import (
	"github.com/zibianqu/eino_study/internal/model"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(doc *model.Document) error
	GetByID(docID string) (*model.Document, error)
	GetByPath(filePath string) (*model.Document, error)
	List(offset, limit int) ([]*model.Document, int64, error)
	Update(doc *model.Document) error
	Delete(docID string) error
	UpdateSyncState(docID string, ragState, entityState int) error
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(doc *model.Document) error {
	return r.db.Create(doc).Error
}

func (r *documentRepository) GetByID(docID string) (*model.Document, error) {
	var doc model.Document
	err := r.db.Where("doc_id = ?", docID).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) GetByPath(filePath string) (*model.Document, error) {
	var doc model.Document
	err := r.db.Where("file_path = ?", filePath).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) List(offset, limit int) ([]*model.Document, int64, error) {
	var docs []*model.Document
	var total int64

	if err := r.db.Model(&model.Document{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("ctime DESC").Find(&docs).Error
	return docs, total, err
}

func (r *documentRepository) Update(doc *model.Document) error {
	return r.db.Save(doc).Error
}

func (r *documentRepository) Delete(docID string) error {
	return r.db.Where("doc_id = ?", docID).Delete(&model.Document{}).Error
}

func (r *documentRepository) UpdateSyncState(docID string, ragState, entityState int) error {
	return r.db.Model(&model.Document{}).Where("doc_id = ?", docID).Updates(map[string]interface{}{
		"sync_rag_state":    ragState,
		"sync_enity_state": entityState,
	}).Error
}