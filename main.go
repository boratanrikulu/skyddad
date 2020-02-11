package main

import (
  "github.com/jinzhu/gorm"
)

var db *gorm.DB

func main() {
  db = Connect()
  Migrate(db)
  defer db.Close()
}
