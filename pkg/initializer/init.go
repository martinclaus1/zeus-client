package initializer

import (
	"errors"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const appName = "zeus-client"

var lock = &sync.Mutex{}

type BaseConfig struct {
	configFilePath string
	logFilePath    string
	machineId      string
}

var singleInstance *BaseConfig

func GetBaseConfig() *BaseConfig {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			var err error
			singleInstance, err = createInstance()
			if err != nil {
				panic(err)
			}
		}
	}

	return singleInstance
}

func createInstance() (*BaseConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not get user home directory. error: %v", err))
	}

	configFolderPath := homeDir + filepath.FromSlash("/."+appName)
	logFolderPath := configFolderPath + filepath.FromSlash("/logs")

	if _, err = os.Stat(logFolderPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(logFolderPath, os.ModePerm)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("could not create config directory. error: %v", err))
		}
	}

	machineId, err := machineid.ProtectedID(appName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not get machine id. error: %v", err))
	}

	if len(machineId) < 24 {
		machineId += strings.Repeat("a", 24-len(machineId))
	} else {
		machineId = machineId[0:24]
	}

	return &BaseConfig{
		configFilePath: configFolderPath + filepath.FromSlash("/config"),
		logFilePath:    logFolderPath + filepath.FromSlash("/zeus-client.log"),
		machineId:      machineId,
	}, nil
}

func (receiver *BaseConfig) ConfigFilePath() string {
	return receiver.configFilePath
}

func (receiver *BaseConfig) LogFilePath() string {
	return receiver.logFilePath
}

func (receiver *BaseConfig) MachineId() string {
	return receiver.machineId
}
