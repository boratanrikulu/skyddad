## Skyddad.

> It means "protected" in Swedish.

This program was written for Cryptology lesson that's given at Pamukkale University.

## Features.

- All mails are encrypted by using [**Stream Cipher Algorithm (CFB)**](https://golang.org/pkg/crypto/cipher/#Stream).  
- You can simply see if mails are changed. Mail hashes are calculated by using SHA-256 algorithm [**crypto/sha256**](https://golang.org/pkg/crypto/sha256).
- Mails are signed by using [**ED25519 Algorithm**](https://golang.org/pkg/crypto/ed25519/). That is an automatic operation. When you signup to system, private and public keys are created. When you send a mail, the mail will be signed by using your private key. Users checks received email's signatures by checking from-user's public key.

## Installation

> This project needs to Go (At least 1.14) to compile.  
  Download from [**here**](https://golang.org/dl/).

**Get the repo.**

```bash
go get github.com/boratanrikulu/skyddad
```

**Set your DB.**

This project needs Postgresql DB.  
You need to create a database named **skyddad**.

**Set your env file.**

You need to set database information to env file.  
Set `.env` file to wherever you use the skyddad command or `${HOME}/.config/skyddad/.env`

There is a env sample: [**here**](/env.sample).

## Usage

```
NAME:
   Skyddad - A mail client that keeps you safe.

USAGE:
   skyddad [global options] command [command options] [arguments...]

COMMANDS:
   mails        Show all mails that were sent by the user.
   send-mail    Send mail to the user.
   sign-up      Sign up to the mail service.
   spam-attack  Attack to the user with spam mails.
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## Testing

#### Creating users.
```bash
skyddad sign-up --username "testing-user-1" --password "user-1-pass"
```

#### Excepted result.
```
(✓) User was created.
  Username: testing-user-1,
  Password: user-1-pass,
```

---

#### Sending mails.
```bash
skyddad send-mail --username "testing-user-1" --password "user-1-pass" \
                  --to-user "testing-user-2" \
                  --body "Top secret message."
```

#### Excepted result.
> Body section would be different.  

```
------------------
(✓) Mail was sent.
	----------
	From: testing-user-1,
	To: testing-user-2
	Date: 2020-04-30 19:34:34.435383788 +0300 +03 m=+0.120306794,
	Hash: f1b8a5f9377b8b77a21eb61234383d5c071aca09cdd20bacbd88dafeef6bf3a4
	Signature: 4b1a105d60dfd23d56f80a0fa298452458e20b1ad3f75b77a6dd6d81f67a44057fce4c8ca8aa6b10af5381aa64297d5190f8ec480dc39b480587751f8617300f
	Body: [ Encrypted ] 5795eb9b062d86a1482aa6c378f5ef1c99e26d98fba30aea99e45ab8fda34ba93c8ca1
------------------
(✓) A mail was sent to "testing-user-1" from "testing-user-2".
```

---

#### Showing e-mails.
```bash
skyddad mails --username "testing-user-2" --password "user-2-pass"
```

#### Excepted result.
```
------------------
To: testing-user-2
	----------
	(✓) Message is not changed. Hash is same.
	(✓) Message is signed by testing-user-1. That's an real signature.
	From: testing-user-1,
	To: testing-user-2
	Date: 2020-04-30 19:34:34.435384 +0300 +03,
	Hash: f1b8a5f9377b8b77a21eb61234383d5c071aca09cdd20bacbd88dafeef6bf3a4
	Signature: 4b1a105d60dfd23d56f80a0fa298452458e20b1ad3f75b77a6dd6d81f67a44057fce4c8ca8aa6b10af5381aa64297d5190f8ec480dc39b480587751f8617300f
	Body: [ Decrypted ] Top secret message.
------------------
(✓) "1" mails are listed for "testing-user-2" user.
```

---

#### Spam attack to the user.
```bash
skyddad spam-attack --username "testing-user-1" --password "user-1-pass" \
                    --to-user "testing-user-2" \
                    --number-of-mails "5"
```

#### Excepted result.
```
------------------
(✓) Mail was sent.
	----------
	From: testing-user-1,
	To: testing-user-2
	Date: 2020-04-30 19:38:13.325474363 +0300 +03 m=+0.117518951,
	Hash: d0c8ac4f1ab67fa73201eef453d51520a72b028c7a23676c6dbeb974ddc14e96
	Signature: 07e0a7b9e57bbc91df760df2cb15ec000cb64cd2e7fb4c15b29698ef7f85b1409b737296f70ae01ef9da9b8955c208dc0567abafd4ac453414071b577559aa05
	Body: [ Encrypted ] 15d0593106daf17ab7093356702ad26f8d939a21f97c7d9ff839d8299aec56945988419bcb156496f09efab36ff002aba8b5715f791a4991660facd1e02d2c95bd82cca6253ef9ef2398906c8739392c52bc32b59dfa11607ccdc97276d5c52ed28785f1198f94edce21eb1bcc39933344789c4edbed6c22f9e346228b2b2a1e494b
	Body Text: Consider Consider uninstallation Our good of a it. I before my original completed from So device have its not 404 
------------------
(✓) Mail was sent.
	----------
	From: testing-user-1,
	To: testing-user-2
	Date: 2020-04-30 19:38:13.330720123 +0300 +03 m=+0.122764712,
	Hash: f70c90e43f6945eb422267ad03d95b5fc1123f2c595b2e50da15f2dd199293f1
	Signature: 6050cd086513f693d69f76b856911d3ecbbf1b4c52f1c00e3b4302ac5e731c482d96c2dd9b558817d03f3e8b7fac90f194c6dce708460db57359cc1094b37c03
	Body: [ Encrypted ] a794b5942125968dd57108849aea3a8ceab839f7c17865216b776762879902f4f161ae3f6397608035514f3976ed07bb26eaa3fb6c77326238b01752b14a3a9b0fb547b2792d50fa9e9018e17f7cad1fc8e8253fe85d7981b76632ea125ce0a8c7fd43f337d4bcc4f5329e0d8567050c79f49bd8ee497e87969e9b9340f537559d53caa72a098c68
	Body Text: the wish now original process, did You, my uninstallation Our found. luck Consider to years found. had been broken back 
------------------
(✓) Mail was sent.
	----------
	From: testing-user-1,
	To: testing-user-2
	Date: 2020-04-30 19:38:13.335537344 +0300 +03 m=+0.127581931,
	Hash: 62201221fd8c5937288a4cbef185b1f1b4f1adb0b5174035a3301d5e610e1979
	Signature: 2cf85768d17647635fb94ed5a12502200536570d6647650d4aa3d00e29f4c4eb6e7c4aa7c454b2828d655ef3faea19f1e4f3090d062b0d44dcd4b8253ec9df08
	Body: [ Encrypted ] 38190a5e927e6c4d3ec391ae6aae54cfc41a9debe228e08444d21498b3b1f6e860e5aa1a53ac1df009ac83de89f6440689e6b83126ddb711126abce5f3d343f40f93cc5aeb73c0f9bc3163313a0ebe1145f79622b7e3d75b4d0331c8936083164cc88f1d21dd744ecef831ebf63cb7e405e0b9dab332701ae4007d3711fc3e250365ba88e997c961
	Body Text: now a its So completed unregistered Consider trusted from I task before back relationship, my settings You, I to before 
------------------
(✓) Mail was sent.
	----------
	From: testing-user-1,
	To: testing-user-2
	Date: 2020-04-30 19:38:13.339615542 +0300 +03 m=+0.131660132,
	Hash: e55b243174be368b35231915f2ef37ed6812c47fb217df549c494ad54022b049
	Signature: c643a34a87c94b738796d767f36c6aaf8d95dc3a45d244d5005d1f0586aea47806d80d5a6ebd43378325d45e43ffafbd1f867267c122dd7edfe4e7be5253150c
	Body: [ Encrypted ] c88593acb64e7f57173df8c480311c360d83c590bd997305cf852bb067cfce33b90122164ada90420ae7f93c9732288cec3be77a4e6ff8ddc992e6ff91a178e87c24d9928127666b0be67b15340b5afe6870b309f4aef77f8039074e75c0362a4eca10b6ebbb7bc66ad4a19be27b33ce1e15a4ac93d453332c9876ac55269c9ecc5acdde9e3cf3ff9de9873c6c
	Body Text: forbidden have now wish luck initially visit, ended. I your to You your have forbidden forbidden had luck spamming initially 
------------------
(✓) Mail was sent.
	----------
	From: testing-user-1,
	To: testing-user-2
	Date: 2020-04-30 19:38:13.343571884 +0300 +03 m=+0.135616472,
	Hash: a5402043cd8b0d4cf772e18dd30c9bcb4637f28b80d28a28ef6ec34b8291bd42
	Signature: 3604a0f11094920356ee88308253852ec9b6d825a464931bff4ca0e2ace2e27835dab3013b58bcc494e8042e11e25ca3a16f9284eff6b9f1070c5113122b3e0f
	Body: [ Encrypted ] 8f5804b26215c9af6841c4f7aba08b4029bf79b194e4c0cedb5449c1d13ccfdf3595f331fb69a53000365b3f6243153a1ae56accab3eb757ba6eb3fa30c86cf7835de8ddbd3a5e572c341a59c7da08ab1394a5da14861a149772612797aae0e2032aefc9338f59d7213eef61e4fb1fe8c420067e1f027b192470d0f94a
	Body Text: 404 Our before in had the its years have files. had five five download in a uninstallation have already five 
------------------
(✓) Spam attack has been completed. "5" mails was sent to "testing-user-2".
```

## To-Do

- [x] Add end-to-end encryption between users.  
- [x] Add spam attack feature. (--spam-attack)  
- [x] Add hash control feature for checking if message is changed.  
- [x] Add electronic signatures for e-mails.
- [ ] Add encryption for user passwords.
