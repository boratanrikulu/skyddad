package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/boratanrikulu/skyddad/controller"
	"github.com/boratanrikulu/skyddad/driver"
	"github.com/boratanrikulu/skyddad/model"
)

// DB variable is exported to use on the whole project.
// Connection is set by using driver/connection.
var (
	DB = driver.Connect()
)

func init() {
	DB = driver.Connect()
}

func main() {
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
							fmt.Println("------------------")
							showRecivedMail(&mail)
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
					if isEmpty(c.String("secret-message")) != isEmpty(c.String("image-path")) || isEmpty(c.String("secret-message")) != isEmpty(c.String("passphrase")) {
						log.Fatal("You need to set secret-message, passphrase and image-path both!")
					}

					currentUser := controller.LogIn(c.String("username"), c.String("password"))
					toUser := model.User{}
					DB.Where("username = ?", c.String("to-user")).First(&toUser)
					if currentUser.Username != "" {
						if toUser.Username != "" {
							result, mail := controller.SendMail(currentUser, toUser, c.String("body"), c.String("key"), c.String("secret-message"), c.String("passphrase"), c.String("watermark"), c.String("image-path"))
							if result {
								fmt.Printf("------------------\n")
								fmt.Println("(✓) Mail was sent.")
								fmt.Printf("\t----------\n")
								setEncryptionInfo(&mail, "[ Encrypted ] ")
								showSendMail(&mail)
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
					// TODO: remove key option
					&cli.StringFlag{
						Name:     "key, k",
						Usage:    "Custom key to encrypt mail. Automatically create key if you do not set custom key.",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "secret-message, sm",
						Usage:    "Secret message to encode into the image.",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "passphrase, pp",
						Usage:    "Passphrase to decode the secret message.",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "watermark, w",
						Usage:    "Text to add image as a watermark.",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "image-path, ip",
						Usage:    "Image path to send to the user.",
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
								result, mail := controller.SendMail(currentUser, toUser, body, c.String("key"), "", "", "", "")
								if result {
									fmt.Printf("------------------\n")
									fmt.Println("(✓) Mail was sent.")
									setEncryptionInfo(&mail, "[ Encrypted ] ")
									fmt.Printf("\t----------\n")
									showSendMail(&mail)
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
			{
				Name:  "set-2fa",
				Usage: "Sets 2fa for your account.",
				Action: func(c *cli.Context) error {
					currentUser := controller.LogIn(c.String("username"), c.String("password"))
					if currentUser.Username != "" {
						// If 2fa is already active, ask for make inactive.
						if currentUser.Is2faActive {
							result := controller.Set2faInactive(&currentUser)
							if result {
								fmt.Println("2FA is inactivated successfully.")
								return nil
							}
							fmt.Println("Operation is canceled. Your account's 2FA is still active.")
						} else {
							// If 2fa is not active, ask for make active.
							result := controller.Set2faActive(&currentUser)
							if result {
								fmt.Println("2FA is activated successfully.")
								return nil
							}
							log.Fatal("Error occur. We can not active 2FA for your account")
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
				},
			},
			{
				Name:  "secret-message",
				Usage: "Sets 2fa for your account.",
				Action: func(c *cli.Context) error {
					imagePath := c.String("image-path")
					passphrase := c.String("passphrase")
					secretImage, err := ioutil.ReadFile(imagePath)
					if err != nil {
						log.Fatalf("Error occur while reading image: %v", err)
					}
					secretMessage := controller.GetSecretMessageFromImage(secretImage, passphrase)
					fmt.Printf("Secret message: %v\n", secretMessage)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "image-path, i",
						Usage:    "Image path to decode secret message.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "passphrase, p",
						Usage:    "Passphrase to read the secret message.",
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

func showRecivedMail(mail *model.Mail) {
	controller.SetMailUser(mail)
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
	setEncryptionInfo(mail, "[ Decrypted ] ")

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

	if mail.IsContainImage {
		file := controller.CreateTempFile(mail.Image)
		path, err := filepath.Abs(file.Name())
		if err == nil {
			fmt.Println("\t----------")
			fmt.Println("\tImage: It containes an secret image.")
			fmt.Printf("\tImage saved at: \"%v\"\n", path)
			fmt.Println("\t----------")
		}
	}
}

func showSendMail(mail *model.Mail) {
	controller.SetMailUser(mail)
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

	if mail.IsContainImage {
		fmt.Println("\t----------")
		fmt.Println("\tImage: Secret image is attach to mail.")
		fmt.Println("\t----------")
	}
}

func isEmpty(s string) bool {
	if len(strings.TrimSpace(s)) == 0 {
		return true
	}
	return false
}
