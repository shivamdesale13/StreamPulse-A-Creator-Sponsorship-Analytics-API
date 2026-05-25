package queue

import "context"

type EventType string

const (
	EventDealCreated   EventType = "deal.created"
	EventDealActivated EventType = "deal.activated"
	EventDealCompleted EventType = "deal.completed"
	EventDealCancelled EventType = "deal.cancelled"
	EventDealUpdated   EventType = "deal.updated"
)

type DealEvent struct {
	Type      EventType `json:"type"`
	DealID    string    `json:"deal_id"`
	CreatorID string    `json:"creator_id"`
	SponsorID string    `json:"sponsor_id"`
	Value     float64   `json:"value"`
	Timestamp string    `json:"timestamp"`
}

// Queue is the abstraction over SQS (prod) or an in-memory channel (dev/test).
type Queue interface {
	Publish(ctx context.Context, event DealEvent) error
	Subscribe() <-chan DealEvent
	Close()
}
