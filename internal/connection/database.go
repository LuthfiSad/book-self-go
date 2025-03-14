package connection

import (
	"database/sql"
	"fmt"
	"go-rest-api/domain"
	"go-rest-api/internal/config"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabase(conf config.Database) (*sql.DB, *gorm.DB) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		conf.Host,
		conf.User,
		conf.Pass,
		conf.Name,
		conf.Port,
		conf.Tz,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed open connection to db: ", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("failed open connection to db: ", err.Error())
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to create gorm db: ", err.Error())
	}

	autoMigrate(gormDB)

	return db, gormDB
}

func autoMigrate(DB *gorm.DB) {
	err := DB.AutoMigrate(&domain.User{}, &domain.Book{}, &domain.BookStock{}, &domain.Media{}, &domain.Journal{}, &domain.Charge{}, &domain.Customer{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("✅ Database migrated successfully!")
}
