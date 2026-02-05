package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/eino_study/internal/app/service"
	"github.com/zibianqu/eino_study/pkg/api"
)

type DocumentHandler struct {
	docService service.DocumentService
}

func NewDocumentHandler(docService service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		docService: docService,
	}
}

// Upload handles document upload
func (h *DocumentHandler) Upload(c *gin.Context) {
	var req api.DocumentUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	doc, err := h.docService.UploadDocument(req.FilePath, req.DocName)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, doc)
}

// Get handles get document by ID
func (h *DocumentHandler) Get(c *gin.Context) {
	docID := c.Param("id")
	if docID == "" {
		BadRequest(c, "document id is required")
		return
	}

	doc, err := h.docService.GetDocument(docID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}

	Success(c, doc)
}

// List handles list documents
func (h *DocumentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	docs, total, err := h.docService.ListDocuments(page, perPage)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessWithPage(c, docs, total, page, perPage)
}

// Delete handles delete document
func (h *DocumentHandler) Delete(c *gin.Context) {
	docID := c.Param("id")
	if docID == "" {
		BadRequest(c, "document id is required")
		return
	}

	if err := h.docService.DeleteDocument(docID); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{"message": "document deleted successfully"})
}

// Process handles document processing (RAG sync)
func (h *DocumentHandler) Process(c *gin.Context) {
	docID := c.Param("id")
	if docID == "" {
		BadRequest(c, "document id is required")
		return
	}

	if err := h.docService.ProcessDocument(docID); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, gin.H{"message": "document processing started"})
}