package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/api/internal/models"
	"github.com/streampulse/api/internal/service"
)

type SponsorHandler struct {
	svc *service.SponsorService
}

func NewSponsorHandler(svc *service.SponsorService) *SponsorHandler {
	return &SponsorHandler{svc: svc}
}

func (h *SponsorHandler) Create(c *gin.Context) {
	var req models.CreateSponsorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sponsor, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sponsor)
}

func (h *SponsorHandler) List(c *gin.Context) {
	sponsors, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sponsors": sponsors, "count": len(sponsors)})
}

func (h *SponsorHandler) GetByID(c *gin.Context) {
	sponsor, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sponsor)
}

func (h *SponsorHandler) Update(c *gin.Context) {
	var req models.UpdateSponsorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sponsor, err := h.svc.Update(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sponsor)
}

func (h *SponsorHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
