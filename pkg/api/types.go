package api

// Response represents a standard API response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse represents a paginated response
type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}

// DocumentUploadRequest represents document upload request
type DocumentUploadRequest struct {
	FilePath string `json:"file_path" binding:"required"`
	DocName  string `json:"doc_name"`
}

// QueryRequest represents a query request
type QueryRequest struct {
	Query   string `json:"query" binding:"required"`
	TopK    int    `json:"top_k,omitempty"`
	Stream  bool   `json:"stream,omitempty"`
}

// QueryResponse represents a query response
type QueryResponse struct {
	Answer      string              `json:"answer"`
	Sources     []SourceInfo        `json:"sources,omitempty"`
	Usage       *UsageInfo          `json:"usage,omitempty"`
}

// SourceInfo represents source document info
type SourceInfo struct {
	DocID      string  `json:"doc_id"`
	DocName    string  `json:"doc_name"`
	Content    string  `json:"content"`
	Similarity float64 `json:"similarity,omitempty"`
}

// UsageInfo represents token usage info
type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ListDocumentsRequest represents list documents request
type ListDocumentsRequest struct {
	Page    int `form:"page" binding:"min=1"`
	PerPage int `form:"per_page" binding:"min=1,max=100"`
}