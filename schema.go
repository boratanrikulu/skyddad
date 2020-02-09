package main

import (
  "github.com/jinzhu/gorm"
)

type User struct {
  gorm.Model
  Username string `gorm:"unique;unique_index;not null"`
  Password string `gorm:"not null"`
  Mails []Mail `gorm:"foreignkey:ToRefer;association_foreignkey:ID"`
}

type Mail struct {
  gorm.Model
  From User `gorm:"association_foreignkey:FromRefer"`
  FromRefer string
  To User `gorm:"association_foreignkey:ToRefer"`
  ToRefer string
}

func Migrate(db *gorm.DB) {
  db.AutoMigrate(&User{}, &Mail{})
}
