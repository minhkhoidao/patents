package models

type Result struct {
	ID                             uint      `gorm:"primaryKey"`
	PatentID                       uint      `gorm:"index"`
	ChallengeEmbedding             []float64 `gorm:"type:float[]"`
	SolutionEmbedding              []float64 `gorm:"type:float[]"`
	ExampleEmbedding               []float64 `gorm:"type:float[]"`
	PreviousWorkChallengeEmbedding []float64 `gorm:"type:float[]"`
}

func (Result) TableName() string {
	return "public.results"
}
