package pipeline

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streampulse/api/internal/models"
	"github.com/streampulse/api/internal/queue"
	"github.com/streampulse/api/internal/repository"
)

// DealTracker consumes deal events from the queue and drives async side-effects
// (analytics seeding, audit logging) using a fixed goroutine pool.
type DealTracker struct {
	q          queue.Queue
	analytics  repository.AnalyticsRepository
	deals      repository.DealRepository
	workerPool int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewDealTracker(
	q queue.Queue,
	analytics repository.AnalyticsRepository,
	deals repository.DealRepository,
	workerPool int,
) *DealTracker {
	ctx, cancel := context.WithCancel(context.Background())
	return &DealTracker{
		q:          q,
		analytics:  analytics,
		deals:      deals,
		workerPool: workerPool,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (dt *DealTracker) Start() {
	events := dt.q.Subscribe()
	for i := 0; i < dt.workerPool; i++ {
		dt.wg.Add(1)
		go dt.worker(events)
	}
	log.Printf("[pipeline] deal tracker started with %d workers", dt.workerPool)
}

func (dt *DealTracker) Stop() {
	dt.cancel()
	dt.wg.Wait()
	log.Println("[pipeline] deal tracker stopped")
}

func (dt *DealTracker) worker(events <-chan queue.DealEvent) {
	defer dt.wg.Done()
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}
			dt.processEvent(event)
		case <-dt.ctx.Done():
			return
		}
	}
}

func (dt *DealTracker) processEvent(event queue.DealEvent) {
	ctx := context.Background()
	switch event.Type {
	case queue.EventDealCreated:
		log.Printf("[pipeline] deal %s created — value $%.2f", event.DealID, event.Value)

	case queue.EventDealActivated:
		// Seed an empty analytics record so the deal shows up in summaries immediately.
		a := &models.Analytics{
			ID:         uuid.New().String(),
			DealID:     event.DealID,
			RecordedAt: time.Now(),
		}
		if err := dt.analytics.Create(ctx, a); err != nil {
			log.Printf("[pipeline] failed to seed analytics for deal %s: %v", event.DealID, err)
		} else {
			log.Printf("[pipeline] analytics seeded for activated deal %s", event.DealID)
		}

	case queue.EventDealCompleted:
		log.Printf("[pipeline] deal %s completed — final value $%.2f", event.DealID, event.Value)

	case queue.EventDealCancelled:
		log.Printf("[pipeline] deal %s cancelled", event.DealID)
	}
}
