package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/api/internal/models"
	"github.com/streampulse/api/internal/service"
)

type CreatorHandler struct {
	svc       *service.CreatorService
	analytics *service.AnalyticsService
}

func NewCreatorHandler(svc *service.CreatorService, analytics *service.AnalyticsService) *CreatorHandler {
	return &CreatorHandler{svc: svc, analytics: analytics}
}

func (h *CreatorHandler) Create(c *gin.Context) {
	var req models.CreateCreatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creator, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, creator)
}

func (h *CreatorHandler) List(c *gin.Context) {
	creators, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"creators": creators, "count": len(creators)})
}

func (h *CreatorHandler) GetByID(c *gin.Context) {
	creator, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, creator)
}

func (h *CreatorHandler) Update(c *gin.Context) {
	var req models.UpdateCreatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creator, err := h.svc.Update(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, creator)
}

func (h *CreatorHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CreatorHandler) GetAnalytics(c *gin.Context) {
	entries, err := h.analytics.GetCreatorAnalytics(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"analytics": entries, "count": len(entries)})
}
