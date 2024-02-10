package env

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Frank-Mayer/env/internal"
)

func Import(url string, key string) error {
	// http get url
	res, err := http.Get(url)
	if err != nil {
		return errors.Join(fmt.Errorf("could not get url %s", url), err)
	}
	defer res.Body.Close()

	// read body
	enc, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Join(fmt.Errorf("could not read body"), err)
	}

	// decode key
	binKey, err := internal.Atob(key)
	if err != nil {
		return errors.Join(fmt.Errorf("could not decode key"), err)
	}

	// check if the key is valid
	if len(binKey) != 32 {
		return fmt.Errorf("invalid key length %d, expected 32", len(binKey))
	}
	b64Key := internal.Btoa(binKey)
	if b64Key != key {
		return fmt.Errorf("invalid key format")
	}

	// decrypt
	dec, err := internal.Decrypt(binKey, enc)
	if err != nil {
		return errors.Join(fmt.Errorf("could not decrypt"), err)
	}

	// stingify
	decStr := string(dec)

	// parse env
	lines := strings.Split(decStr, "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		i := strings.Index(line, "=")
		if i == -1 {
			return fmt.Errorf("invalid line %s", line)
		}
		key := line[:i]
		value := line[i+1:]
		err := os.Setenv(key, value)
		if err != nil {
			return errors.Join(fmt.Errorf("could not set env %s", key), err)
		}
	}

	return nil
}
