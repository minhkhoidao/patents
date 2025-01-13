package main

import (
	"flag"
	"fmt"
	"go-nginx/command"
	"go-nginx/controllers"
	"go-nginx/database"
	"go-nginx/repositories"
	"go-nginx/routes"
	"go-nginx/services"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting application...")

	// Add logging for database connection parameters
	fmt.Printf("Patent DB Connection params: host=%s dbname=%s port=%s\n",
		os.Getenv("PATENT_DB_HOST"),
		os.Getenv("PATENT_DB_NAME"),
		os.Getenv("PATENT_DB_PORT"),
	)

	patentDBURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("PATENT_DB_HOST"),     // Use env var instead of hardcoded "patent_db"
		os.Getenv("PATENT_DB_USER"),     // Use env var instead of "postgres"
		os.Getenv("PATENT_DB_PASSWORD"), // Use env var instead of "123456"
		os.Getenv("PATENT_DB_NAME"),     // Use env var instead of "patent_db"
		os.Getenv("PATENT_DB_PORT"),     // Use env var instead of "5432"
	)

	// Add better error handling for patent DB
	patentDB, err := gorm.Open(postgres.Open(patentDBURL), &gorm.Config{})
	if err != nil {
		fmt.Printf("Patent DB connection error: %v\n", err)
		panic(fmt.Sprintf("failed to connect patent database: %v", err))
	}

	dataDir := flag.String("data", "", "Directory containing the patent CSV file")
	flag.Parse()

	// Only attempt import if dataDir is provided
	if *dataDir != "" {
		if _, err := os.Stat(*dataDir); os.IsNotExist(err) {
			log.Fatalf("Error: Directory %s does not exist", *dataDir)
		}

		cmd := command.NewImportPatentCommand(patentDB)
		cmd.SetDataDir(*dataDir)

		if err := cmd.Execute(); err != nil {
			log.Fatalf("Failed to import patents: %v", err)
		}

		log.Println("Patent import completed successfully")
	}

	err = database.Migrate(patentDB, nil)
	if err != nil {
		panic("failed to migrate database")
	}

	patentRepo := repositories.NewPatentRepository(patentDB)
	patentService := services.NewPatentService(patentRepo)
	patentController := controllers.NewPatentController(patentService)
	router := routes.SetupRouter(patentController)

	router.Run(":8080") // Gin's way to start the server
}
