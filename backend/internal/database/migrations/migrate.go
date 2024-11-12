// internal/database/migrations/migrate.go
package migrations

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func RunMigrations(gormDB *gorm.DB) error {
	log.Println("Running migrations...")

	// Get the underlying *sql.DB from GORM
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying *sql.DB: %v", err)
	}

	// Define the migrations directory
	migrationsPath := "file://internal/database/migrations/files"

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath, // Use the relative path
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
