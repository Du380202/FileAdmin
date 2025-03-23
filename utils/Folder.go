package utils

import (
	"os"
)

func CheckFolder(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return false, err
		}
	}
	return true, nil
}
