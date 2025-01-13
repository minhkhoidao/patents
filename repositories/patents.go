package repositories

import (
	"context"
	"fmt"
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

	// Get the query execution plan
	var result []map[string]interface{}
	r.db.WithContext(ctx).
		Raw("EXPLAIN ANALYZE SELECT * FROM patents WHERE file_date >= ? AND file_date <= ?",
			start_date, end_date).
		Scan(&result)

	// Print the execution plan
	for _, row := range result {
		fmt.Printf("%v\n", row)
	}

	// Execute the actual query
	err := r.db.WithContext(ctx).Where("file_date >= ? AND file_date <= ?", start_date, end_date).Find(&patents).Error
	return patents, err
}
