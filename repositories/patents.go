package repositories

import (
	"context"
	"go-nginx/models"

	"gorm.io/gorm"
)

type PatentRepository struct {
	db *gorm.DB
}

func NewPatentRepository(db *gorm.DB) *PatentRepository {
	return &PatentRepository{db: db}
}

func (r *PatentRepository) GetDataGraph(ctx context.Context, start_date string, end_date string) ([]models.Patent, error) {
	var patents []models.Patent
	err := r.db.WithContext(ctx).Where("file_date >= ? AND file_date <= ?", start_date, end_date).Find(&patents).Error
	return patents, err
}
