package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// type Config struct {
// 	// If you want to configure YAML based envs
// 	err error
// }

func MustLoad() error {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "local"
	}
	envConfigPath := filepath.Join("internal", "envs", fmt.Sprintf(".env.%s", env))
	err := godotenv.Load(envConfigPath)

	// return &Config{
	// 	err: err,
	// }
	return err
}
