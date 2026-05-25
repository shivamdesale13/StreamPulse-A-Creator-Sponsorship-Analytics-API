package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streampulse/api/internal/models"
	"github.com/streampulse/api/internal/service"
)

type DealHandler struct {
	svc       *service.DealService
	analytics *service.AnalyticsService
}

func NewDealHandler(svc *service.DealService, analytics *service.AnalyticsService) *DealHandler {
	return &DealHandler{svc: svc, analytics: analytics}
}

func (h *DealHandler) Create(c *gin.Context) {
	var req models.CreateDealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deal, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, deal)
}

func (h *DealHandler) List(c *gin.Context) {
	deals, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deals": deals, "count": len(deals)})
}

func (h *DealHandler) GetByID(c *gin.Context) {
	deal, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deal)
}

func (h *DealHandler) UpdateStatus(c *gin.Context) {
	var req models.UpdateDealStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deal, err := h.svc.UpdateStatus(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deal)
}

func (h *DealHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DealHandler) GetAnalytics(c *gin.Context) {
	entries, err := h.analytics.GetForDeal(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"analytics": entries, "count": len(entries)})
}

func (h *DealHandler) RecordAnalytics(c *gin.Context) {
	var req models.CreateAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entry, err := h.analytics.RecordForDeal(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entry)
}
