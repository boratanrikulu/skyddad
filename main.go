package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/urfave/cli/v2"

	"github.com/boratanrikulu/skyddad/controller"
	"github.com/boratanrikulu/skyddad/driver"
	"github.com/boratanrikulu/skyddad/model"
)

var DB *gorm.DB

func main() {
	DB = driver.Connect()
	model.Migrate(DB)
	defer DB.Close()

	app := &cli.App{
		Name:  "Skyddad",
		Usage: "A mail client that keeps you safe.",
		Commands: []*cli.Command{
			{
				Name:  "mails",
				Usage: "Show all mails that were sent by the user.",
				Action: func(c *cli.Context) error {
					currentUser := controller.LogIn(c.String("username"), c.String("password"))
					if currentUser.Username != "" {
						fmt.Printf("------------------\n")
						fmt.Printf("To: %v\n", currentUser.Username)
						take := c.String("take")
						if take == "" {
							take = "-1"
						}
						mails := controller.Mails(currentUser, take)
						for _, mail := range mails {
							controller.SetMailUser(&mail)
							fmt.Printf("\t----------\n")
							if controller.IsChanged(mail.Body, mail.Hash) {
								fmt.Println("\t(!) Message is changed. Hash is NOT same!")
							} else {
								fmt.Println("\t(✓) Message is not changed. Hash is same.")
							}
							if len(mail.Signature) != 0 {
								if controller.IsSignatureReal(mail.From.PublicKey, []byte(mail.Hash), mail.Signature) {
									fmt.Printf("\t(✓) Message is signed by %v. That's an real signature.\n", mail.From.Username)
								} else {
									fmt.Printf("\t(!) Message is signed by %v. But that signature is FAKE!\n", mail.From.Username)
								}
							} else {
								fmt.Println("Message is not signed.")
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
					currentUser := controller.LogIn(c.String("username"), c.String("password"))
					toUser := model.User{}
					DB.Where("username = ?", c.String("to-user")).First(&toUser)
					if currentUser.Username != "" {
						if toUser.Username != "" {
							result, mail := controller.SendMail(currentUser, toUser, c.String("body"), c.String("key"))
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
					result, user := controller.SingUp(c.String("username"), c.String("password"))
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
					currentUser := controller.LogIn(c.String("username"), c.String("password"))
					toUser := model.User{}
					DB.Where("username = ?", c.String("to-user")).First(&toUser)
					if currentUser.Username != "" {
						if toUser.Username != "" {
							count, err := strconv.Atoi(c.String("number-of-mails"))
							if err != nil {
								log.Fatal("Number of emails value must be an integer.")
							}
							for i := 0; i < count; i++ {
								body := controller.RandomMails()
								result, mail := controller.SendMail(currentUser, toUser, body, c.String("key"))
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

func setEncryptionInfo(mail *model.Mail, info string) {
	body := ""
	if mail.IsEncrypted == true {
		body += info
	} else {
		body += "[ Plain Text ] "
	}
	body += mail.Body
	mail.Body = body
}

func showUser(user model.User) {
	DB.First(&user)
	fmt.Printf("\tUsername: %v,\n\tPassword: %v,\n", user.Username, user.Password)
}

func showMail(mail model.Mail) {
	controller.SetMailUser(&mail)
	signature := "\n\tSignature: There is no signature for this mail."
	if len(mail.Signature) != 0 {
		signature = fmt.Sprintf("\n\tSignature: %x", mail.Signature)
	}
	fmt.Printf("\tFrom: %v,\n\tTo: %v\n\tDate: %v,\n\tHash: %v%v\n\tBody: %v\n",
		mail.From.Username,
		mail.To.Username,
		mail.CreatedAt,
		mail.Hash,
		signature,
		mail.Body)
}
