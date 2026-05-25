package models

import "time"

type Sponsor struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Industry     string    `json:"industry"`
	Website      string    `json:"website"`
	ContactEmail string    `json:"contact_email"`
	Budget       float64   `json:"budget"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateSponsorRequest struct {
	Name         string  `json:"name"          binding:"required"`
	Industry     string  `json:"industry"      binding:"required"`
	Website      string  `json:"website"`
	ContactEmail string  `json:"contact_email" binding:"required,email"`
	Budget       float64 `json:"budget"        binding:"required,gt=0"`
}

type UpdateSponsorRequest struct {
	Name    string  `json:"name"`
	Website string  `json:"website"`
	Budget  float64 `json:"budget"`
}
