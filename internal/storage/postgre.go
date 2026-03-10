package storage

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConnection(host, user, password, dbname, port string) (*gorm.DB, error) {
	var db *gorm.DB 
	var err error 
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	for i := 0; i <= 4; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
		if err == nil {
			return db, nil
		}
		time.Sleep(1 * time.Second)
	}
	return nil, err 
}