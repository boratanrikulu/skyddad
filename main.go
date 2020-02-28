package main

import (
  "fmt"
  "os"
  "log"
  "github.com/jinzhu/gorm"
  "github.com/urfave/cli/v2"
)

var db *gorm.DB

func main() {
  db = Connect()
  Migrate(db)
  defer db.Close()

  app := &cli.App{
    Name: "Skyddad",
    Usage: "A mail client that keeps you safe.",
    Commands: []*cli.Command{
      {
        Name:    "mails",
        Usage:   "Show all mails that were sent by the user.",
        Action:  func(c *cli.Context) error {
          currentUser := LogIn(c.String("username"), c.String("password"));
          if currentUser.Username != "" {
            fmt.Printf("To: %v\n", currentUser.Username)
            take := c.String("take")
            if take == "" {
              take = "-1"
            }
            showMails(Mails(currentUser, take))
          } else {
            fmt.Println("(!) Incorrect username or password.")
          }
          return nil
        },
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:  "username, u",
            Usage: "Your username to use mail service.",
            Required: true,
          },
          &cli.StringFlag{
            Name:  "password, p",
            Usage: "Your password to use mail service.",
            Required: true,
          },
          &cli.StringFlag{
            Name: "take, t",
            Usage: "Mail limit to take.",
          },
        },
      },
      {
        Name:    "send-mail",
        Usage:   "Send mail to the user.",
        Action:  func(c *cli.Context) error {
          currentUser := LogIn(c.String("username"), c.String("password"));
          toUser := User{}
          db.Where("username = ?", c.String("to-user")).First(&toUser)
          if currentUser.Username != "" {
            if toUser.Username != "" {
              result, mail := SendMail(currentUser, toUser, c.String("body"))
              if result {
                fmt.Println("(✓) Mail was sent.")
                showMail(mail)
              } else {
                fmt.Println("(!) Error occur while sending mail.")
              }
            } else {
              fmt.Printf("(!) There is no user to send mail: %v.\n", toUser.Username)
            }
          } else {
            fmt.Println("(!) Incorrect username or password.")
          }
          return nil
        },
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:  "username, u",
            Usage: "Your username to use mail service.",
            Required: true,
          },
          &cli.StringFlag{
            Name:  "password, p",
            Usage: "Your password to use mail service.",
            Required: true,
          },
          &cli.StringFlag{
            Name:  "to-user, t",
            Usage: "Username to send mail.",
            Required: true,
          },
          &cli.StringFlag{
            Name: "body, b",
            Usage: "Body for the mail.",
            Required: true,
          },
        },
      },
      {
        Name: "sign-up",
        Usage: "Sign up to the mail service.",
        Action: func(c *cli.Context) error {
          result, user := SingUp(c.String("username"), c.String("password"))
          if result {
            fmt.Println("(✓) User was created.")
            showUser(user)
          } else {
            fmt.Println("(!) Error occur while creating user.")
          }
          return nil
        },
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:  "username, u",
            Usage: "Your username to use mail service.",
            Required: true,
          },
          &cli.StringFlag{
            Name:  "password, p",
            Usage: "Your password to use mail service.",
            Required: true,
          },
        },
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}

func showUser(user User) {
  db.First(&user)
  fmt.Printf("\tUsername: %v,\n\tPassword: %v,\n", user.Username, user.Password)
}

func showMail(mail Mail) {
  from := User{}
  db.Model(&mail).Association("From").Find(&from)
  body := ""
  if mail.IsEncrypted == true {
    body += "[ Encrypted Text ] "
  } else {
    body += "[ Plain Text ] "
  }
  body += mail.Body
  fmt.Printf("\tFrom: %v,\n\tDate: %v,\n\tBody: %v\n", from.Username, mail.CreatedAt, body)
}

func showMails(mails []Mail) {
  for i, mail := range mails {
    if i != 0 {
      fmt.Printf("\t----------\n")
    }
    showMail(mail)
  }
}
