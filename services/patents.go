package services

import (
	"context"
	"fmt"
	"go-nginx/models"
	"go-nginx/repositories"
	"time"
)

type PatentService struct {
	patentRepo   *repositories.PatentRepository
	cacheService *RedisService
}

func NewPatentService(repo *repositories.PatentRepository, cache *RedisService) *PatentService {
	return &PatentService{
		patentRepo:   repo,
		cacheService: cache,
	}
}

func (s *PatentService) GetDataGraph(ctx context.Context, startDate, endDate string) ([]models.Patent, error) {
	cacheKey := fmt.Sprintf("patents:graph:%s:%s", startDate, endDate)

	// Try to get from cache
	var patents []models.Patent
	err := s.cacheService.Get(ctx, cacheKey, &patents)
	if err == nil {
		return patents, nil
	}

	// If not in cache, get from repository
	patents, err = s.patentRepo.GetDataGraph(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if len(patents) > 0 {
		err = s.cacheService.Set(ctx, cacheKey, patents, time.Hour)
		if err != nil {
			// Log error but don't fail the request
			// log.Printf("Failed to cache patents: %v", err)
		}
	}

	return patents, nil
}
