package main

func LogIn(username string, password string) bool {
  user := User {}
  db.Where("username = ? AND password = ?", username, password).First(&user)
  if user.Username == "" {
    return false
  }
  return true
}

func SingUp(username string, password string) bool {
  user := User{
    Username: username,
    Password: password,
  }
  db.Create(&user)
  return !db.NewRecord(user)
}

func SendMail(from User, to User, body string) bool {
  mail := Mail {
    From: from,
    To: to,
    Body: body,
  }
  db.Create(&mail)
  return !db.NewRecord(mail)
}

func UserMails(user User) []Mail {
  mails := []Mail{}
  db.Model(&user).Association("Mails").Find(&mails)
  return mails
}
