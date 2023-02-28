package configs

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"healthcheck/pkg/models"
	"io"
	"os"
	"path/filepath"
)

func InitConfig(path string) error {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("couldn't get mainPath")
		return nil
	}
	/// TODO for Windows
	err = godotenv.Load(dir + path)
	if err != nil {
		logrus.Errorf("failed godotenv.Load: %v", err)
		return nil
	}

	return nil
}

func ReadUrlConfigs(path string) (*models.Configuration, error) {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("failed os.Getwd: %s", err)
		return nil, err
	}
	jsonFile, err := os.Open(filepath.Join(dir, path))
	if err != nil {
		logrus.Errorf("failed os.Open: %s", err)
		return nil, err
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		logrus.Errorf("failed ioutil.ReadAll: %s", err)
		return nil, err
	}
	var configs models.Configuration

	err = json.Unmarshal(byteValue, &configs)
	if err != nil {
		logrus.Errorf("failed json.Unmarshal: %s", err)
		return nil, err
	}
	return &configs, nil
}
