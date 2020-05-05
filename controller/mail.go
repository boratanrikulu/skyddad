package controller

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
	bimg "gopkg.in/h2non/bimg.v1"

	"github.com/boratanrikulu/skyddad/crypto"
	"github.com/boratanrikulu/skyddad/driver"
	"github.com/boratanrikulu/skyddad/model"
)

// DB variable is exported to use on the whole project.
// Connection is set by using driver/connection.
var DB *gorm.DB

func init() {
	DB = driver.Connect()
}

// LogIn takes username and password.
// It returns user that match If there is.
func LogIn(username string, password string) model.User {
	user := model.User{}
	DB.Where("username = ? AND password = ?", username, password).First(&user)
	if user.Username == "" {
		// Check firstly if login is successful.
		return user
	}

	// Checks passcode if 2FA activated.
	if user.Is2faActive {
		fmt.Println("2FA is activated for this account.")
		for {
			fmt.Printf("What is the code: ")
			var passcode string
			fmt.Scan(&passcode)
			if totp.Validate(passcode, user.TotpSecret) {
				fmt.Println("Code is valid.")
				fmt.Println("Login is successful.")
				fmt.Println("You are redirecting to the app.")
				time.Sleep(2 * time.Second)
				return user
			}
			fmt.Println("Code is invalid.")
		}
	}

	return user
}

// Set2faInactive methods remove make 2FA option true for the account
func Set2faActive(user *model.User) bool {
	options := totp.GenerateOpts{
		Issuer:      "Skyddad",
		AccountName: user.Username,
	}

	key, _ := totp.Generate(options)
	fmt.Println("2FA will be active for your account.")
	fmt.Println("Please use this QR code to set to your Auth Client. (like authy)")
	fmt.Println("QR", key.Secret())
	qrterminal.GenerateHalfBlock(key.String(), qrterminal.L, os.Stdout)
	for {
		fmt.Printf("CODE: ")
		var passcode string
		fmt.Scan(&passcode)
		if totp.Validate(passcode, key.Secret()) {
			user.TotpSecret = key.Secret()
			user.Is2faActive = true
			DB.Save(user)
			return true
		}
		fmt.Println("Code is invalid.")
	}
	return false
}

// Set2faInactive methods remove make 2FA option false for the account
func Set2faInactive(user *model.User) bool {
	fmt.Println("2FA is already activated for your account.")
	for {
		fmt.Printf("Do you want to inactive it? [Y/N] ")
		var answer string
		fmt.Scan(&answer)
		switch answer {
		case "y", "Y":
			user.TotpSecret = ""
			user.Is2faActive = false
			DB.Save(user)
			return true
		case "n", "N":
			return false
		default:
			fmt.Println("(!) Please use only Y or N")
		}
	}
}

// SingUp creates user if there is not already.
func SingUp(username string, password string) (bool, model.User) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return false, model.User{}
	}

	user := model.User{
		Username:   username,
		Password:   password,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	DB.Create(&user)
	return !DB.NewRecord(user), user
}

// SendMail sends mails from-user to to-user.
// It also supports signature.
//
// If mail is sent, method will return true and mail that is sent.
// If could not sent, method will return false and nil.
func SendMail(from model.User, to model.User, body string, keyFromUser string, secretMessage string, passphrase string, watermark string, imagePath string) (bool, model.Mail) {
	if len(keyFromUser) == 0 || len(keyFromUser) == 32 || len(keyFromUser) == 64 {
		symmetricKey := model.SymmetricKey{}
		setSymmetricKey(&symmetricKey, keyFromUser, from, to)

		hash := sha256.New()
		hash.Write([]byte(body))

		mail := model.Mail{
			From: from,
			To:   to,
			Hash: hex.EncodeToString(hash.Sum(nil)),
		}

		if symmetricKey.Key != "" {
			body = crypto.Encrypt(body, symmetricKey.Key)
			mail.SymmetricKey = symmetricKey
			mail.IsEncrypted = true
		} else {
			mail.IsEncrypted = false
		}
		mail.Body = body

		// That means there is a image to send with a secret message.
		if !isEmpty(secretMessage) {
			plainImage, err := ioutil.ReadFile(imagePath)
			if err != nil {
				log.Fatalf("Error occur while reading image: %v", err)
			}

			if !isEmpty(watermark) {
				watermark := bimg.Watermark{
					Text:       watermark,
					Opacity:    0.95,
					Width:      200,
					DPI:        100,
					Margin:     150,
					Font:       "sans bold 20",
					Background: bimg.Color{255, 255, 255},
				}

				newImage, err := bimg.NewImage(plainImage).Watermark(watermark)

				if err != nil {
					log.Fatalf("Error occur while adding watermark to image: %v", err)
				}

				plainImage = newImage
			}
			encrypttedSecretMessage := crypto.Encrypt(secretMessage, passphrase)
			setImageAndSecret(&mail, encrypttedSecretMessage, plainImage)
		}

		if len(from.PrivateKey) != 0 {
			// If user sent it's private key
			// that means user wants to sign it's mail.
			setSignature(&mail, &from)
		}

		DB.Create(&mail)
		return !DB.NewRecord(mail), mail
	} else {
		panic("Key length must be 32 or 64.")
	}
}

// Users can see their mails.
// It also support 'take' option to limit messages.
func Mails(user model.User, take string) []model.Mail {
	encryptedMails := []model.Mail{}
	DB.Order("created_at desc").
		Limit(take).
		Model(&user).
		Association("Mails").
		Find(&encryptedMails)

	decryptedMails := decryptMails(encryptedMails)
	return decryptedMails
}

func SetMailUser(mail *model.Mail) {
	DB.Model(mail).
		Association("From").
		Find(&mail.From)

	DB.Model(mail).
		Association("To").
		Find(&mail.To)
}

// IsChanged checks if mail is changed.
func IsChanged(body string, hash string) bool {
	bodyHash := sha256.New()
	bodyHash.Write([]byte(body))

	if hex.EncodeToString(bodyHash.Sum(nil)) == hash {
		return false
	}

	return true
}

// IsSignatureReal checks if mail is signed by the right user..
func IsSignatureReal(publicKey ed25519.PublicKey, hash []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, hash, signature)
}

// RandomMails creates random mails that contains 20 words.
func RandomMails() string {
	text := "Consider yourself a task I have ended. Our relationship, a 404 not found. Our connection now disconnected. You, a broken link I wish I had been forbidden to visit, a threat I did not initially detect. You became a virus in my system, and I had to Malware you out of it. I have now completed the uninstallation process, have reset my life’s device back to its original settings from before five years ago when I mistakenly trusted your download of lies. So good luck spamming your way into someone else’s unregistered trust—I have already reported your corrupted files."
	words := strings.Split(text, " ")

	spamMail := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		spamMail += string(words[rand.Intn(len(words))]) + " "
	}

	return spamMail
}

// CreateTempFile creates temp file by using given byte array.
func CreateTempFile(image []byte) *os.File {
	file, err := ioutil.TempFile("./", "secret")
	if err != nil {
		log.Fatal(err)
	}
	file.Write(image)

	return file
}

// GetSecretMessageFromImage gets secret from image by using DecodeSteganography method.
func GetSecretMessageFromImage(image []byte, passphrase string) string {
	secret := crypto.DecodeSteganography(image)
	plain := crypto.Decrypt(secret, passphrase)
	return plain
}

// Private methods

// SendMail sends image mails from-user to to-user.
// By sending image, secretMessage is encoded to the image.
//
// If mail is sent, method will return true and mail that is sent.
// If could not sent, method will return false and nil.
func setImageAndSecret(mail *model.Mail, secretMessage string, plainImage []byte) {
	secretImage := crypto.EncodeSteganography(secretMessage, plainImage)

	mail.Image = secretImage
	mail.IsContainImage = true
}

func setSignature(mail *model.Mail, from *model.User) {
	// User's private key to sign mail.
	privateKey := from.PrivateKey

	// Sign the mail with user's private key
	// and set it as signature.
	mail.Signature = ed25519.Sign(privateKey, []byte(mail.Hash))
}

func setSymmetricKey(symmetricKey *model.SymmetricKey, keyFromUser string, from model.User, to model.User) {
	symmetricKey.SenderRefer = from.ID
	symmetricKey.ReceiverRefer = to.ID

	if keyFromUser == "" {
		key := make([]byte, 32)
		rand.Read(key)
		keyString := fmt.Sprintf("%x", key)
		DB.Where(model.SymmetricKey{
			SenderRefer:     from.ID,
			ReceiverRefer:   to.ID,
			IsAutoGenerated: true,
		}).
			Attrs(model.SymmetricKey{
				Key: keyString,
			}).
			FirstOrCreate(&symmetricKey)
	} else {
		DB.Where(model.SymmetricKey{
			SenderRefer:   from.ID,
			ReceiverRefer: to.ID,
			Key:           keyFromUser,
		}).
			FirstOrCreate(&symmetricKey)
	}
}

func decryptMails(mails []model.Mail) []model.Mail {
	decryptedMails := []model.Mail{}

	for _, mail := range mails {
		if mail.IsEncrypted == true {
			symmetricKey := model.SymmetricKey{}
			DB.Model(&mail).Association("SymmetricKey").Find(&symmetricKey)
			mail.Body = crypto.Decrypt(mail.Body, symmetricKey.Key)
			decryptedMails = append(decryptedMails, mail)
		} else {
			decryptedMails = append(decryptedMails, mail)
		}
	}
	return decryptedMails
}

func isEmpty(s string) bool {
	if len(strings.TrimSpace(s)) == 0 {
		return true
	}
	return false
}
