package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zibianqu/eino_study/internal/app/service"
	"github.com/zibianqu/eino_study/pkg/api"
)

type QueryHandler struct {
	ragService service.RAGService
}

func NewQueryHandler(ragService service.RAGService) *QueryHandler {
	return &QueryHandler{
		ragService: ragService,
	}
}

// Query handles knowledge base query
func (h *QueryHandler) Query(c *gin.Context) {
	var req api.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	if req.TopK <= 0 {
		req.TopK = 5
	}

	resp, err := h.ragService.Query(req.Query, req.TopK)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, resp)
}