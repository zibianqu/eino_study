package service

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/internal/model"
	"github.com/zibianqu/eino_study/internal/pkg/utils"
	"gorm.io/gorm"
)

type DocumentService interface {
	UploadDocument(filePath, docName string) (*model.Document, error)
	GetDocument(docID string) (*model.Document, error)
	ListDocuments(page, perPage int) ([]*model.Document, int64, error)
	DeleteDocument(docID string) error
	ProcessDocument(docID string) error
}

type documentService struct {
	docRepo   repository.DocumentRepository
	chunkRepo repository.ChunkRepository
	entityRepo repository.EntityRepository
}

func NewDocumentService(
	docRepo repository.DocumentRepository,
	chunkRepo repository.ChunkRepository,
	entityRepo repository.EntityRepository,
) DocumentService {
	return &documentService{
		docRepo:   docRepo,
		chunkRepo: chunkRepo,
		entityRepo: entityRepo,
	}
}

func (s *documentService) UploadDocument(filePath, docName string) (*model.Document, error) {
	// Check if file exists
	if !utils.FileExists(filePath) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// Check if document already exists
	existing, err := s.docRepo.GetByPath(filePath)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("document already exists with path: %s", filePath)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing document: %w", err)
	}

	// Calculate file hash
	fileHash, err := utils.MD5File(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	// Generate document ID
	docID := utils.MD5String(filePath)

	// Set document name
	if docName == "" {
		docName = filepath.Base(filePath)
	}

	// Get file type
	fileType := filepath.Ext(filePath)

	// Create document
	doc := &model.Document{
		DocID:           docID,
		DocName:         docName,
		DocHash:         fileHash,
		FilePath:        filePath,
		FileType:        fileType,
		SyncRagState:    0,
		SyncEntityState: 0,
		CTime:           time.Now(),
	}

	if err := s.docRepo.Create(doc); err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return doc, nil
}

func (s *documentService) GetDocument(docID string) (*model.Document, error) {
	doc, err := s.docRepo.GetByID(docID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return doc, nil
}

func (s *documentService) ListDocuments(page, perPage int) ([]*model.Document, int64, error) {
	offset := (page - 1) * perPage
	docs, total, err := s.docRepo.List(offset, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list documents: %w", err)
	}
	return docs, total, nil
}

func (s *documentService) DeleteDocument(docID string) error {
	// Delete chunks
	if err := s.chunkRepo.DeleteByDocID(docID); err != nil {
		return fmt.Errorf("failed to delete chunks: %w", err)
	}

	// Delete entities
	if err := s.entityRepo.DeleteByDocID(docID); err != nil {
		return fmt.Errorf("failed to delete entities: %w", err)
	}

	// Delete document
	if err := s.docRepo.Delete(docID); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

func (s *documentService) ProcessDocument(docID string) error {
	// TODO: Implement document processing with Eino
	// 1. Load document
	// 2. Split into chunks
	// 3. Generate embeddings
	// 4. Store in vector database
	// 5. Extract entities
	// 6. Update sync state
	return fmt.Errorf("not implemented yet")
}