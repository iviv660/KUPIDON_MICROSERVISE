package handler

import (
	"context"
	"net/http"
	"service3/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	uc     *usecase.Usecase
	router *gin.Engine
}

func NewMatchHandler(uc *usecase.Usecase, router *gin.Engine) *MatchHandler {
	handler := &MatchHandler{uc: uc, router: router}
	router.POST("/like/:id1/:id2", handler.Like)
	return handler
}

func (h *MatchHandler) Like(c *gin.Context) {
	fromUserID, err1 := strconv.ParseInt(c.Param("id1"), 10, 64)
	toUserID, err2 := strconv.ParseInt(c.Param("id2"), 10, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err := h.uc.Like(context.Background(), fromUserID, toUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like saved"})
}
