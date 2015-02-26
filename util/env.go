package util

import (
	"os"
)

func GetOpt(name, def string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return def
}
