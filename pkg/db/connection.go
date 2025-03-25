package db

import (
	"fmt"

	"github.com/joejosephvarghese/message/server/pkg/config"
	"github.com/joejosephvarghese/message/server/pkg/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase initializes and returns a new gorm DB instance
func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		fmt.Println("❌ Database connection failed:", err)
		return nil, err
	}

	// Check if DB connection is working
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("❌ Failed to get database instance:", err)
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		fmt.Println("❌ Database ping failed:", err)
		return nil, err
	}

	fmt.Println("✅ Database connected successfully!")

	// Migrate tables
	err = db.AutoMigrate(
		domain.User{},
		domain.RefreshSession{},
		domain.Chat{},
		domain.Message{},
	)

	if err != nil {
		fmt.Println("❌ Migration failed:", err)
		return nil, err
	}

	fmt.Println("✅ Migration successful!")

	return db, nil
}
