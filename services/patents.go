package services

import (
	"context"
	"go-nginx/models"
	"go-nginx/repositories"
)

type PatentService struct {
	repo *repositories.PatentRepository
}

func NewPatentService(repo *repositories.PatentRepository) *PatentService {
	return &PatentService{repo: repo}
}

func (s *PatentService) GetDataGraph(ctx context.Context, start_date string, end_date string) ([]models.Patent, error) {
	patents, err := s.repo.GetDataGraph(ctx, start_date, end_date)
	if err != nil {
		return nil, err
	}

	return patents, nil
}
