package user

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"github.com/martinclaus1/zeus-client/pkg/initializer"
	log "github.com/sirupsen/logrus"
	"os"
)

type User struct {
	Username string
	Password string
}

func (config *User) Save() {
	content, err := json.Marshal(config)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not marshal config.")
	}

	encryptedContent, err := encrypt(string(content))
	if err != nil {
		log.WithField("error", err).Fatalln("Could not encrypt config.")
	}

	err = os.WriteFile(initializer.GetBaseConfig().ConfigFilePath(), []byte(encryptedContent), 0600)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not write config file.")
	}
}

func ReadConfig() *User {
	content, err := os.ReadFile(initializer.GetBaseConfig().ConfigFilePath())
	if os.IsNotExist(err) {
		return &User{}
	} else if err != nil {
		log.WithField("error", err).Fatalln("Could not read config file.")
	}

	decryptedContent, _ := decrypt(string(content))

	var config User
	err = json.Unmarshal(decryptedContent, &config)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not unmarshal config.")
	}

	return &config
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not decode config.")
	}
	return data
}

func encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(initializer.GetBaseConfig().MachineId()))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

func decrypt(text string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(initializer.GetBaseConfig().MachineId()))
	if err != nil {
		return nil, err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return plainText, nil
}
