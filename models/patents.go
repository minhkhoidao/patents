package models

import "time"

type Patent struct {
	ID                             uint      `gorm:"primaryKey"`
	OriginalID                     int       `gorm:"column:original_id;comment:Original CSV Unnamed: 0 column"`
	PatentName                     string    `gorm:"type:varchar(500);not null;column:patent_name"`
	Author                         string    `gorm:"type:varchar(200)"`
	IPC                            string    `gorm:"type:text;column:ipc"` // International Patent Classification
	FileDate                       time.Time `gorm:"index;column:file_date"`
	RetroDate                      time.Time
	Abstract                       string   `gorm:"type:text;column:abstract"`
	Solution                       string   `gorm:"type:text;column:solution"`
	Challenge                      string   `gorm:"type:text;column:challenge"`
	PreviousWork                   string   `gorm:"type:text;column:previous_work"`
	Example                        string   `gorm:"type:text;column:example"`
	PreviousWorkChallenge          string   `gorm:"type:text;column:previous_work_challenge"`
	ChallengeTsne2dOne             *float64 `gorm:"column:challenge_tsne_2d_one"`
	ChallengeTsne2dTwo             *float64 `gorm:"column:challenge_tsne_2d_two"`
	ExampleTsne2dOne               *float64 `gorm:"column:example_tsne_2d_one"`
	ExampleTsne2dTwo               *float64 `gorm:"column:example_tsne_2d_two"`
	PreviousworkchallengeTsne2dOne *float64 `gorm:"column:previousworkchallenge_tsne_2d_one"`
	PreviousworkchallengeTsne2dTwo *float64 `gorm:"column:previousworkchallenge_tsne_2d_two"`
	SolutionTsne2dOne              *float64 `gorm:"column:solution_tsne_2d_one"`
	SolutionTsne2dTwo              *float64 `gorm:"column:solution_tsne_2d_two"`
}

// TableName specifies the table name for the Patent model
func (Patent) TableName() string {
	return "public.patents"
}
