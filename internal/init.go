package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func Init() (bool, error) {
	config := &Config{
		configSchema(),
		[]Variable{
			{"DB_PORT", "3306"},
		},
		[]Profile{
			{
				"dev",
				[]Variable{
					{"DB_HOST", "localhost"},
					{"DB_USER", "root"},
					{"DB_PASS", "12345"},
				},
			},
			{
				"prod",
				[]Variable{
					{"DB_HOST", "123.456.789.0"},
					{"DB_USER", "root"},
					{"DB_PASS", "ab123!"},
				},
			},
		}}

	fmt.Println("Waiting for user to edit config")
	if err, changed := EditStruct(config); err != nil {
		return false, errors.Join(errors.New("failed to edit config"), err)
	} else if !changed {
		fmt.Println("No changes made to config, aborting")
		return false, nil
	}

	fmt.Println("Initializing config")

	// create api keys for each profile
	profileKeys := make(map[string][]byte, len(config.Profiles))
	for _, p := range config.Profiles {
		// check if profile already has a key
		if _, ok := profileKeys[p.Name]; ok {
			return false, errors.New("Duplicate profile name: " + p.Name)
		}
		// create key
		if k, err := NewPassword(); err != nil {
			return false, errors.Join(errors.New("failed to create new password"), err)
		} else {
			profileKeys[p.Name] = k
		}
	}

	// create master key
	masterKey, err := NewPassword()
	if err != nil {
		return false, errors.Join(errors.New("failed to create new password"), err)
	} else {
		// encrypt main config as json
		fn := "config.json"
		f, err := os.Create(fn)
		if err != nil {
			return false, errors.Join(errors.New("failed to create file"), err)
		}
		defer f.Close()
		jsonStr, err := json.Marshal(config)
		if err != nil {
			return false, errors.Join(errors.New("failed to marshal json"), err)
		}
		encrypted, err := Encrypt(masterKey, jsonStr)
		if err != nil {
			return false, errors.Join(errors.New("failed to encrypt json"), err)
		}
		_, err = f.Write(encrypted)
		if err != nil {
			return false, errors.Join(errors.New("failed to write to file"), err)
		}
	}

	if _, err = os.Stat("profiles"); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("profiles", 0777); err != nil {
				return false, errors.Join(errors.New("failed to create profiles directory"), err)
			}
		} else {
			return false, errors.Join(errors.New("failed to check for profiles directory"), err)
		}
	} else {
		// remove all files in profiles directory
		d, err := os.ReadDir("profiles")
		if err != nil {
			return false, errors.Join(errors.New("failed to read profiles directory"), err)
		}
		if len(d) > 0 {
			fmt.Println("Removing existing profiles")
		}
		for _, de := range d {
			fmt.Printf("Deleting %s\n", de.Name())
			if err := os.RemoveAll(filepath.Join("profiles", de.Name())); err != nil {
				return false, errors.Join(errors.New("failed to remove file"), err)
			}
		}
	}

	fmt.Println("Creating profiles")
	for _, p := range config.Profiles {
		fmt.Printf("Creating profile for %s\n", p.Name)
		if err := func() error {
			// encrypt profile as json
			fn := p.Name + ".env"
			fp := filepath.Join("profiles", fn)
			f, err := os.Create(fp)
			if err != nil {
				return errors.Join(errors.New("failed to create file"), err)
			}
			defer f.Close()
			var sb strings.Builder
			for _, v := range p.Variables {
				sb.WriteString(v.Key + "=" + v.Value + "\n")
			}
			envStr := sb.String()
			encrypted, err := Encrypt(profileKeys[p.Name], []byte(envStr))
			if err != nil {
				return errors.Join(errors.New("failed to encrypt json"), err)
			}
			_, err = f.Write(encrypted)
			if err != nil {
				return errors.Join(errors.New("failed to write to file"), err)
			}
			return nil
		}(); err != nil {
			return false, err
		}
	}

	fmt.Println("Creating main profile containing all profile keys")
	if err := func() error {
		fn := "main.json"
		f, err := os.Create(fn)
		if err != nil {
			return errors.Join(errors.New("failed to create file"), err)
		}
		defer f.Close()
		keys := make(map[string]string, len(profileKeys))
		for k, v := range profileKeys {
			keys[k] = Btoa(v)
		}
		jsonStr, err := json.Marshal(keys)
		if err != nil {
			return errors.Join(errors.New("failed to marshal json"), err)
		}
		encrypted, err := Encrypt(masterKey, []byte(jsonStr))
		if err != nil {
			return errors.Join(errors.New("failed to encrypt json"), err)
		}
		_, err = f.Write(encrypted)
		if err != nil {
			return errors.Join(errors.New("failed to write to file"), err)
		}
		return nil
	}(); err != nil {
		return false, err
	}

	fmt.Print("Done\n\n")

	// ask if master key should be printed or copied to clipboard
	switch UserChoise("What do you want to do with the master key?", "Print to console", "Copy to clipboard") {
	case 1:
		fmt.Println("Master key:", Btoa(masterKey))
	case 2:
		if err := clipboard.WriteAll(Btoa(masterKey)); err != nil {
			return false, errors.Join(errors.New("failed to copy to clipboard"), err)
		}
	default:
		return false, errors.New("Error in user input")
	}

	return true, nil
}
