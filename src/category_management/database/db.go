package database

import (
	"log"
	"os"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/category_management/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Migrator interface {
	AutoMigrate(dst ...interface{}) error
}

var DB *gorm.DB

// StartDB initializes the database connection and applies any pending migrations.
// It retrieves the database DSN from the environment variables,
// opens a connection to the MySQL database, and runs automatic migrations
// to ensure that the database schema is up to date.
func StartDB() {
	dsn := os.Getenv("DB_DSN")
	db := OpenDBConnection(dsn, mysql.Open)
	MakeMigrations(db)
	DB = db
}

// openDBConnection opens a connection to a MySQL database using the provided DSN.
// It will panic if an error occurs while opening the connection.
func OpenDBConnection(dsn string, opener func(dsn string)(gorm.Dialector)) *gorm.DB {
	db, err := gorm.Open(opener(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("Connected to the database")
	return db
}

// makeMigrations performs automatic database schema migration for all models defined in the application.
// It ensures that the database schema is up to date with the latest model definitions.
func MakeMigrations(db Migrator) {
	err := db.AutoMigrate(&models.Category{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	log.Println("Migrations done")
}
