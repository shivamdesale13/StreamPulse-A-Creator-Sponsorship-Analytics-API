package models

import "time"

type DealStatus string

const (
	DealStatusPending   DealStatus = "pending"
	DealStatusActive    DealStatus = "active"
	DealStatusCompleted DealStatus = "completed"
	DealStatusCancelled DealStatus = "cancelled"
)

type Deal struct {
	ID          string     `json:"id"`
	CreatorID   string     `json:"creator_id"`
	SponsorID   string     `json:"sponsor_id"`
	Status      DealStatus `json:"status"`
	Value       float64    `json:"value"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateDealRequest struct {
	CreatorID   string    `json:"creator_id"  binding:"required"`
	SponsorID   string    `json:"sponsor_id"  binding:"required"`
	Value       float64   `json:"value"       binding:"required,gt=0"`
	StartDate   time.Time `json:"start_date"  binding:"required"`
	EndDate     time.Time `json:"end_date"    binding:"required"`
	Description string    `json:"description"`
}

type UpdateDealStatusRequest struct {
	Status DealStatus `json:"status" binding:"required"`
}
