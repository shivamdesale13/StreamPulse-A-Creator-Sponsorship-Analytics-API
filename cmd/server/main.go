package main

import (
	"log"

	"github.com/streampulse/api/internal/auth"
	"github.com/streampulse/api/internal/config"
	"github.com/streampulse/api/internal/handlers"
	"github.com/streampulse/api/internal/pipeline"
	memqueue "github.com/streampulse/api/internal/queue/memory"
	memrepo "github.com/streampulse/api/internal/repository/memory"
	"github.com/streampulse/api/internal/router"
	"github.com/streampulse/api/internal/service"
)

func main() {
	cfg := config.Load()

	// Repositories (in-memory; swap for DynamoDB clients in prod)
	creatorRepo := memrepo.NewCreatorRepository()
	sponsorRepo := memrepo.NewSponsorRepository()
	dealRepo := memrepo.NewDealRepository()
	analyticsRepo := memrepo.NewAnalyticsRepository()
	userRepo := memrepo.NewUserRepository()

	// Queue (channel-based; swap for SQS client in prod)
	q := memqueue.NewQueue(cfg.QueueBuffer)

	// Pipeline — goroutine workers consume deal events from the queue
	tracker := pipeline.NewDealTracker(q, analyticsRepo, dealRepo, cfg.WorkerPoolSize)
	tracker.Start()
	defer tracker.Stop()

	// Auth
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Services
	authSvc := service.NewAuthService(userRepo, jwtManager)
	creatorSvc := service.NewCreatorService(creatorRepo)
	sponsorSvc := service.NewSponsorService(sponsorRepo)
	dealSvc := service.NewDealService(dealRepo, q)
	analyticsSvc := service.NewAnalyticsService(analyticsRepo, dealRepo)

	// Handlers
	authH := handlers.NewAuthHandler(authSvc)
	creatorH := handlers.NewCreatorHandler(creatorSvc, analyticsSvc)
	sponsorH := handlers.NewSponsorHandler(sponsorSvc)
	dealH := handlers.NewDealHandler(dealSvc, analyticsSvc)
	analyticsH := handlers.NewAnalyticsHandler(analyticsSvc)

	r := router.New(jwtManager, authH, creatorH, sponsorH, dealH, analyticsH)

	log.Printf("StreamPulse API listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
