package models

import "time"

type Platform string

const (
	PlatformYouTube   Platform = "youtube"
	PlatformTwitch    Platform = "twitch"
	PlatformInstagram Platform = "instagram"
	PlatformTikTok    Platform = "tiktok"
)

type Creator struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Platform        Platform  `json:"platform"`
	ChannelURL      string    `json:"channel_url"`
	SubscriberCount int64     `json:"subscriber_count"`
	Category        string    `json:"category"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateCreatorRequest struct {
	Name            string   `json:"name"             binding:"required"`
	Platform        Platform `json:"platform"         binding:"required"`
	ChannelURL      string   `json:"channel_url"      binding:"required"`
	SubscriberCount int64    `json:"subscriber_count"`
	Category        string   `json:"category"         binding:"required"`
	Email           string   `json:"email"            binding:"required,email"`
}

type UpdateCreatorRequest struct {
	Name            string `json:"name"`
	ChannelURL      string `json:"channel_url"`
	SubscriberCount int64  `json:"subscriber_count"`
	Category        string `json:"category"`
}
