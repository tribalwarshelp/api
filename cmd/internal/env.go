package internal

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"strings"
)

func LoadENVFiles() error {
	for _, filename := range [...]string{".env.local", ".env"} {
		if err := godotenv.Load(filename); err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			return errors.Wrap(err, "godotenv.Load")
		}
	}

	return nil
}
