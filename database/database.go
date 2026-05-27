package database

import (
 "fmt"

 "gorm.io/driver/postgres"
 "gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
 dsn := "host=localhost user=postgres password=WoBaSeMiMa1 dbname=postgres port=5432 sslmode=disable"
 var err error
 DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
 if err != nil {
  panic("Failed to connect to database!")
 }
 fmt.Println("Connected to the database!")

}