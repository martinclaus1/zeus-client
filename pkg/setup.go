package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/denisbrodbeck/machineid"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

var configFilePath = ""
var machineId = ""
var configFolderPath = ""
var logFolderPath = ""

const appName = "zeus-client"
const configFileName = "config"

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.WithField("error", err).Fatalln("Could not get user home directory.")
	}

	configFolderPath = homeDir + filepath.FromSlash("/."+appName)
	logFolderPath = configFolderPath + filepath.FromSlash("/logs")
	if _, err = os.Stat(logFolderPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(logFolderPath, os.ModePerm)
		if err != nil {
			log.WithField("error", err).Fatalln("Could not create config directory.")
		}
	}
	configFilePath = configFolderPath + filepath.FromSlash("/"+configFileName)

	machineId, err = machineid.ProtectedID(appName)
	if len(machineId) < 24 {
		machineId += strings.Repeat("a", 24-len(machineId))
	} else {
		machineId = machineId[0:24]
	}

	if err != nil {
		log.WithField("error", err).Fatalln("Could not get machine id.")
	}
}

func (config *Config) Save() {
	content, err := json.Marshal(config)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not marshal config.")
	}

	encryptedContent, err := encrypt(string(content))
	if err != nil {
		log.WithField("error", err).Fatalln("Could not encrypt config.")
	}

	err = os.WriteFile(configFilePath, []byte(encryptedContent), 0600)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not write config file.")
	}
}

func ReadConfig() *Config {
	content, err := os.ReadFile(configFilePath)
	if os.IsNotExist(err) {
		return &Config{}
	} else if err != nil {
		log.WithField("error", err).Fatalln("Could not read config file.")
	}

	decryptedContent, _ := decrypt(string(content))

	var config Config
	err = json.Unmarshal(decryptedContent, &config)
	if err != nil {
		log.WithField("error", err).Fatalln("Could not unmarshal config.")
	}

	return &config
}

func GetLogfilePath() string {
	return logFolderPath + filepath.FromSlash("/zeus-client.log")
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
	block, err := aes.NewCipher([]byte(machineId))
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
	block, err := aes.NewCipher([]byte(machineId))
	if err != nil {
		return nil, err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return plainText, nil
}

type Config struct {
	Username string
	Password string
}
