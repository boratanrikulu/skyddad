package main

import (
  "fmt"
  "os"
  "log"

  "github.com/jinzhu/gorm"
  "github.com/joho/godotenv"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

func loadEnv() {
  err := godotenv.Load()
  if err != nil {
    err = godotenv.Load(os.Getenv("HOME") + "/.config/skyddad/.env")
    if err != nil {
      log.Fatal("Error loading .env file")
    }
  }
}

func Connect() *gorm.DB {
  loadEnv()
  dbinfo := fmt.Sprintf("host=%v port=%v user=%v password=%v, dbname=%v sslmode=%v",
                        os.Getenv("DB_HOST"),
                        os.Getenv("DB_PORT"),
                        os.Getenv("DB_USER"),
                        os.Getenv("DB_PASSWORD"),
                        os.Getenv("DB_NAME"),
                        os.Getenv("DB_SSLMODE"))
  db, err := gorm.Open("postgres", dbinfo)
  if err != nil {
    panic(err.Error())
  }
  return db
}
