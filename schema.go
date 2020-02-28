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
  From User `gorm:"foreignkey:FromRefer"`
  FromRefer uint `gorm:"not null"`
  To User `gorm:"foreignkey:ToRefer"`
  ToRefer uint `gorm:"not null"`
  Body string `gorm:"not null"`
  IsEncrypted bool `gorm:"not null;default:false"`
}

type SymmetricKey struct {
  gorm.Model
  Sender User `gorm:"foreignkey:SenderRefer;"`
  SenderRefer uint `gorm:"not null"`
  Receiver User `gorm:"foreignkey:ReceiverRefer;"`
  ReceiverRefer uint `gorm:"not null"`
  Key string `gorm:"not null"`
}

func Migrate(db *gorm.DB) {
  db.AutoMigrate(&User{}, &Mail{}, &SymmetricKey{})
}
