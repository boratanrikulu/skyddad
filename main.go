package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/urfave/cli/v2"
)

var db *gorm.DB

func main() {
	db = Connect()
	Migrate(db)
	defer db.Close()

	app := &cli.App{
		Name:  "Skyddad",
		Usage: "A mail client that keeps you safe.",
		Commands: []*cli.Command{
			{
				Name:  "mails",
				Usage: "Show all mails that were sent by the user.",
				Action: func(c *cli.Context) error {
					currentUser := LogIn(c.String("username"), c.String("password"))
					if currentUser.Username != "" {
						fmt.Printf("------------------\n")
						fmt.Printf("To: %v\n", currentUser.Username)
						take := c.String("take")
						if take == "" {
							take = "-1"
						}
						mails := Mails(currentUser, take)
						for _, mail := range mails {
							fmt.Printf("\t----------\n")
							if isChanged(mail.Body, mail.Hash) {
								fmt.Println("\t(!) Message is changed.")
							} else {
								fmt.Println("\t(✓) Message is not changed.")
							}
							setEncryptionInfo(&mail, "[ Decrypted ] ")
							showMail(mail)
						}
						fmt.Printf("------------------\n")
						fmt.Printf("(✓) \"%v\" mails are listed for \"%v\" user.\n", len(mails), c.String("username"))
					} else {
						fmt.Println("(!) Incorrect username or password.")
					}
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "username, u",
						Usage:    "Your username to use mail service.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password, p",
						Usage:    "Your password to use mail service.",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "take, t",
						Usage: "Mail limit to take.",
					},
				},
			},
			{
				Name:  "send-mail",
				Usage: "Send mail to the user.",
				Action: func(c *cli.Context) error {
					currentUser := LogIn(c.String("username"), c.String("password"))
					toUser := User{}
					db.Where("username = ?", c.String("to-user")).First(&toUser)
					if currentUser.Username != "" {
						if toUser.Username != "" {
							result, mail := SendMail(currentUser, toUser, c.String("body"), c.String("key"))
							if result {
								fmt.Printf("------------------\n")
								fmt.Println("(✓) Mail was sent.")
								fmt.Printf("\t----------\n")
								setEncryptionInfo(&mail, "[ Encrypted ] ")
								showMail(mail)
								fmt.Printf("------------------\n")
								fmt.Printf("(✓) A mail was sent to \"%v\" from \"%v\".\n", currentUser.Username, toUser.Username)
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
						Name:     "username, u",
						Usage:    "Your username to use mail service.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password, p",
						Usage:    "Your password to use mail service.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "to-user, t",
						Usage:    "Username to send mail.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "body, b",
						Usage:    "Body for the mail.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "key, k",
						Usage:    "Custom key to encrypt mail. Automatically create key if you do not set custom key.",
						Required: false,
					},
				},
			},
			{
				Name:  "sign-up",
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
						Name:     "username, u",
						Usage:    "Your username to use mail service.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password, p",
						Usage:    "Your password to use mail service.",
						Required: true,
					},
				},
			},
			{
				Name:  "spam-attack",
				Usage: "Attack to the user with spam mails.",
				Action: func(c *cli.Context) error {
					currentUser := LogIn(c.String("username"), c.String("password"))
					toUser := User{}
					db.Where("username = ?", c.String("to-user")).First(&toUser)
					if currentUser.Username != "" {
						if toUser.Username != "" {
							count, err := strconv.Atoi(c.String("number-of-mails"))
							if err != nil {
								log.Fatal("Number of emails value must be an integer.")
							}
							for i := 0; i < count; i++ {
								body := randomMails()
								result, mail := SendMail(currentUser, toUser, body, c.String("key"))
								if result {
									fmt.Printf("------------------\n")
									fmt.Println("(✓) Mail was sent.")
									setEncryptionInfo(&mail, "[ Encrypted ] ")
									fmt.Printf("\t----------\n")
									showMail(mail)
									fmt.Printf("\tBody Text: %v\n", body)
								} else {
									fmt.Println("(!) Error occur while sending mail.")
								}
							}
							fmt.Printf("------------------\n")
							fmt.Printf("(✓) Spam attack has been completed. \"%v\" mails was sent to \"%v\".\n", count, c.String("to-user"))
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
						Name:     "username, u",
						Usage:    "Your username to use mail service.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password, p",
						Usage:    "Your password to use mail service.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "to-user, t",
						Usage:    "Username to send mail.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "number-of-mails, n",
						Usage:    "How many emails to send.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "key, k",
						Usage:    "Custom key to encrypt mail. Automatically create key if you do not set custom key.",
						Required: false,
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

func setEncryptionInfo(mail *Mail, info string) {
	body := ""
	if mail.IsEncrypted == true {
		body += info
	} else {
		body += "[ Plain Text ] "
	}
	body += mail.Body
	mail.Body = body
}

func showUser(user User) {
	db.First(&user)
	fmt.Printf("\tUsername: %v,\n\tPassword: %v,\n", user.Username, user.Password)
}

func showMail(mail Mail) {
	from := User{}
	to := User{}
	db.Model(&mail).Association("From").Find(&from)
	db.Model(&mail).Association("To").Find(&to)
	fmt.Printf("\tFrom: %v,\n\tTo: %v\n\tDate: %v,\n\tHash: %v\n\tBody: %v\n",
		from.Username,
		to.Username,
		mail.CreatedAt,
		mail.Hash,
		mail.Body)
}
