package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/eino_study/pkg/api"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, api.Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithPage(c *gin.Context, data interface{}, total int64, page, perPage int) {
	c.JSON(http.StatusOK, api.PageResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, api.Response{
		Code:    -1,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}