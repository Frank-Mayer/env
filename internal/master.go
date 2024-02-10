package internal

import (
	"errors"
	"os"
	"path/filepath"
)

func ensureFile() (string, error) {
	if loc, err := filepath.Abs("master.env"); err != nil {
		return "", errors.Join(errors.New("failed to get absolute path for master.env"), err)
	} else {
		// does file exist?
		_, err := os.Stat(loc)
		if os.IsNotExist(err) {
			// create file
			f, err := os.Create(loc)
			if err != nil {
				return "", errors.Join(errors.New("failed to create master.env"), err)
			}
			defer f.Close()
		} else if err != nil {
			return "", errors.Join(errors.New("failed to check for master.env"), err)
		}
		return loc, nil
	}
}
