package database

import (
	"fmt"
	"log"

	"github.com/g-ton/stori-candidate/env"

	"github.com/g-ton/stori-candidate/internal/core/entities"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type service struct {
	conn *gorm.DB
}

// Get connection
func connectDB(c env.EnvApp) *gorm.DB {

	//	dsn := "host=localhost user=postgres password=kuna1234 dbname=db_sinbu_exp port=5432 sslmode=disable TimeZone=America/Mexico_City"
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		c.DB_HOST,
		c.DB_USERNAME,
		c.DB_PASSWORD,
		c.DB_NAME,
		c.DB_PORT,
		c.DB_SSLMODE,
		c.DB_TIMEZONE,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database")
		return nil
	}

	return db
}

func New(ec env.EnvApp) *service {
	return &service{
		conn: connectDB(ec),
	}
}

func (db *service) CreateTransaction(transaction *entities.Transaction) error {
	// To do
	fmt.Println(transaction)
	return nil
}
