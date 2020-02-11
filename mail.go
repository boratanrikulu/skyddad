package main

func LogIn(username string, password string) User {
  user := User {}
  db.Where("username = ? AND password = ?", username, password).First(&user)
  return user
}

func SingUp(username string, password string) bool {
  user := User{
    Username: username,
    Password: password,
  }
  db.Create(&user)
  return !db.NewRecord(user)
}

func SendMail(from User, to User, body string) (bool, Mail) {
  mail := Mail {
    From: from,
    To: to,
    Body: body,
  }
  db.Create(&mail)
  return !db.NewRecord(mail), mail
}

func UserMails(user User, take string) []Mail {
  mails := []Mail{}
  db.Order("created_at desc").Limit(take).Model(&user).Association("Mails").Find(&mails)
  return mails
}
