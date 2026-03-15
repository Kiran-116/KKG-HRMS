package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"hrms/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Init initializes the database connection
func Init() error {
	var err error
	cfg := config.AppConfig.Database

	DB, err = sql.Open("postgres", cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.MaxIdleConns)
	DB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	// Run migrations
	if err = RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// RunMigrations runs all migration files in order
func RunMigrations() error {
	migrationsPath := "./migrations"
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// Try alternative path
		migrationsPath = "../migrations"
		if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
			return fmt.Errorf("migrations directory not found")
		}
	}

	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort migration files
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Run each migration
	for _, fileName := range migrationFiles {
		if err := runMigration(filepath.Join(migrationsPath, fileName)); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", fileName, err)
		}
	}

	log.Printf("Successfully ran %d migrations", len(migrationFiles))
	return nil
}

func createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := DB.Exec(query)
	return err
}

func runMigration(filePath string) error {
	fileName := filepath.Base(filePath)
	version := strings.TrimSuffix(fileName, ".sql")

	// Check if migration already applied
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Migration %s already applied, skipping", fileName)
		return nil
	}

	// Read migration file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Applied migration: %s", fileName)
	return nil
}
