package config

import (
	"github.com/joho/godotenv"
)

func NewEnv() {
	godotenv.Load()
}
