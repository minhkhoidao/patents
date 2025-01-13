package command

import (
	"encoding/csv"
	"fmt"
	"go-nginx/models"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// ImportPatentCommand represents the command for importing patent data
type ImportPatentCommand struct {
	DB      *gorm.DB
	DataDir string
}

// NewImportPatentCommand creates a new import patent command
func NewImportPatentCommand(db *gorm.DB) *ImportPatentCommand {
	return &ImportPatentCommand{
		DB: db,
	}
}

// SetDataDir sets the data directory path
func (cmd *ImportPatentCommand) SetDataDir(path string) {
	cmd.DataDir = path
}

func emptyToNone(s string) *float64 {
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}

// Validate checks if required files exist
func (cmd *ImportPatentCommand) Validate() error {
	if cmd.DataDir == "" {
		return fmt.Errorf("data directory path is required")
	}

	// Only check for CSV file
	csvPath := filepath.Join(cmd.DataDir, "df_clean_dr.csv")
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		return fmt.Errorf("required file not found: df_clean_dr.csv")
	}

	return nil
}

func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02", // YYYY-MM-DD
		"2006/01/02", // YYYY/MM/DD
	}

	var parseErr error
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		} else {
			parseErr = err
		}
	}
	return time.Time{}, parseErr
}

// Execute runs the import command
func (cmd *ImportPatentCommand) Execute() error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Clear existing data
	if err := cmd.DB.Exec("DELETE FROM patents").Error; err != nil {
		return fmt.Errorf("failed to clear patents table: %v", err)
	}
	if err := cmd.DB.Exec("DELETE FROM results").Error; err != nil {
		return fmt.Errorf("failed to clear results table: %v", err)
	}

	// Open and read CSV file
	file, err := os.Open(filepath.Join(cmd.DataDir, "df_clean_dr.csv"))
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %v", err)
	}

	// Create header map for easy lookup
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[h] = i
	}

	// Start a transaction for batch processing
	return cmd.DB.Transaction(func(tx *gorm.DB) error {
		batchSize := 100
		var patents []models.Patent
		var results []models.Result

		id := 1
		for {
			record, err := reader.Read()
			if err != nil {
				break // End of file
			}

			// Parse dates
			var fileDate, retroDate time.Time
			if record[headerMap["FileDate"]] != "" {
				fileDate, err = parseDate(record[headerMap["FileDate"]])
				if err != nil {
					log.Printf("Warning: Failed to parse FileDate '%s': %v", record[headerMap["FileDate"]], err)
				}
			}
			if record[headerMap["RetroDate"]] != "" {
				retroDate, err = parseDate(record[headerMap["RetroDate"]])
				if err != nil {
					log.Printf("Warning: Failed to parse RetroDate '%s': %v", record[headerMap["RetroDate"]], err)
				}
			}

			// Create Patent
			patent := models.Patent{
				ID:                             uint(id),
				PatentName:                     record[headerMap["PatentName"]],
				Author:                         record[headerMap["Author"]],
				IPC:                            record[headerMap["IPC"]],
				FileDate:                       fileDate,
				RetroDate:                      retroDate,
				Abstract:                       record[headerMap["Abstract"]],
				Solution:                       record[headerMap["Solution"]],
				Challenge:                      record[headerMap["Challenge"]],
				PreviousWork:                   record[headerMap["PreviousWork"]],
				Example:                        record[headerMap["Example"]],
				PreviousWorkChallenge:          record[headerMap["PreviousWorkChallenge"]],
				ChallengeTsne2dOne:             emptyToNone(record[headerMap["challenge_tsne-2d-one"]]),
				ChallengeTsne2dTwo:             emptyToNone(record[headerMap["challenge_tsne-2d-two"]]),
				ExampleTsne2dOne:               emptyToNone(record[headerMap["example_tsne-2d-one"]]),
				ExampleTsne2dTwo:               emptyToNone(record[headerMap["example_tsne-2d-two"]]),
				PreviousworkchallengeTsne2dOne: emptyToNone(record[headerMap["previousworkchallenge_tsne-2d-one"]]),
				PreviousworkchallengeTsne2dTwo: emptyToNone(record[headerMap["previousworkchallenge_tsne-2d-two"]]),
				SolutionTsne2dOne:              emptyToNone(record[headerMap["solution_tsne-2d-one"]]),
				SolutionTsne2dTwo:              emptyToNone(record[headerMap["solution_tsne-2d-two"]]),
			}
			patents = append(patents, patent)

			// Create Result
			result := models.Result{
				ID:       uint(id),
				PatentID: uint(id),
			}
			results = append(results, result)

			// Batch insert when reaching batch size
			if len(patents) >= batchSize {
				if err := tx.Create(&patents).Error; err != nil {
					return fmt.Errorf("failed to create patent records: %v", err)
				}
				if err := tx.Create(&results).Error; err != nil {
					return fmt.Errorf("failed to create result records: %v", err)
				}

				log.Printf("Processed %d records", id)
				patents = patents[:0] // Clear the slice while keeping capacity
				results = results[:0]
			}

			id++
		}

		// Insert remaining records
		if len(patents) > 0 {
			if err := tx.Create(&patents).Error; err != nil {
				return fmt.Errorf("failed to create final patent records: %v", err)
			}
			if err := tx.Create(&results).Error; err != nil {
				return fmt.Errorf("failed to create final result records: %v", err)
			}
		}

		return nil
	})
}
