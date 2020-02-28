package main

import (
	"crypto/rand"
	"fmt"
)

func LogIn(username string, password string) User {
  user := User {}
  db.Where("username = ? AND password = ?", username, password).First(&user)
  return user
}

func SingUp(username string, password string) (bool, User) {
  user := User{
    Username: username,
    Password: password,
  }
  db.Create(&user)
  return !db.NewRecord(user), user
}

func SendMail(from User, to User, body string) (bool, Mail) {
  symmetricKey := SymmetricKey{}
	key := make([]byte, 32)
	rand.Read(key)
	key_string := fmt.Sprintf("%x", key)
  db.Where(SymmetricKey{
            SenderRefer: from.ID,
            ReceiverRefer: to.ID,
          }).
     Attrs(SymmetricKey{
            Key: key_string,
          }).
     FirstOrCreate(&symmetricKey)

  if symmetricKey.Key != "" {
    encryptedBody := encrypt(body, symmetricKey.Key)

    mail := Mail {
      From: from,
      To: to,
      Body: encryptedBody,
      IsEncrypted: true,
    }

    db.Create(&mail)
    return !db.NewRecord(mail), mail
  } else {
    return false, Mail{}
  }
}

func Mails(user User, take string) []Mail {
  encryptedMails := []Mail{}
  db.Order("created_at desc").
     Limit(take).
     Model(&user).
     Association("Mails").
     Find(&encryptedMails)

  decryptedMails := decryptMails(encryptedMails)
  return decryptedMails
}

func decryptMails(mails []Mail) []Mail {
  decryptedMails := []Mail{}

  for _, mail := range mails {
    if mail.IsEncrypted == true {
      symmetricKey := SymmetricKey{}
      db.Where(SymmetricKey{
                SenderRefer: mail.FromRefer,
                ReceiverRefer: mail.ToRefer,
              }).
        Find(&symmetricKey)
        mail.Body = decrypt(mail.Body, symmetricKey.Key)
        decryptedMails = append(decryptedMails, mail)
    }
    decryptedMails = append(decryptedMails, mail)
  }
  return decryptedMails
}
