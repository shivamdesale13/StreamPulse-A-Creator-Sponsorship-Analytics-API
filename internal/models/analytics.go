package models

import "time"

type Analytics struct {
	ID             string    `json:"id"`
	DealID         string    `json:"deal_id"`
	Views          int64     `json:"views"`
	Clicks         int64     `json:"clicks"`
	ConversionRate float64   `json:"conversion_rate"`
	Revenue        float64   `json:"revenue"`
	RecordedAt     time.Time `json:"recorded_at"`
}

type CreateAnalyticsRequest struct {
	Views   int64   `json:"views"   binding:"gte=0"`
	Clicks  int64   `json:"clicks"  binding:"gte=0"`
	Revenue float64 `json:"revenue" binding:"gte=0"`
}

type AnalyticsSummary struct {
	TotalViews     int64   `json:"total_views"`
	TotalClicks    int64   `json:"total_clicks"`
	TotalRevenue   float64 `json:"total_revenue"`
	AvgConversion  float64 `json:"avg_conversion_rate"`
	ActiveDeals    int     `json:"active_deals"`
	CompletedDeals int     `json:"completed_deals"`
	TotalDeals     int     `json:"total_deals"`
}
