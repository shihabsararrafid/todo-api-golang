package databases

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
    "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSlMode  string
}

func NewDBConnection(config Config) (*sql.DB, error) {
	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSlMode)
	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database:%w", err)
	}
	log.Println("Successfully connected to database")
	return db, nil
}
func RunMigrations(db *sql.DB, migrationPath string) error {
	absPath, err := filepath.Abs(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
log.Printf(absPath)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	sourceURL := fmt.Sprintf("file://%s", absPath)
     m,err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	 )
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}
	
	if err == migrate.ErrNilVersion {
		log.Println("🔄 No migrations to run - database is empty")
	} else {
		log.Printf("🚀 Database migrations completed successfully - Version: %d, Dirty: %v", version, dirty)
	}
	
	return nil
}
