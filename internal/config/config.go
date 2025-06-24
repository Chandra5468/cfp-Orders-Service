package config

import (
	"fmt"
	"os"

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
	envConfigPath := fmt.Sprintf("internal\\envs\\.env.%s", env) // If you arerunning on ubuntu do sepcify accordingly

	err := godotenv.Load(envConfigPath)

	// return &Config{
	// 	err: err,
	// }
	return err
}
