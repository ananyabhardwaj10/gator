package config
import (
	"os"
	"encoding/json"
	"path/filepath"
)

type Config struct {
	DBURL string           `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err 
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, err 
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config 
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err 
	}

	return cfg, nil 
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	return write(*c)
} 

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err 
	}

	return filepath.Join(home, configFileName), nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err 
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err 
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(cfg); err != nil {
		return err 
	}

	return nil 
}