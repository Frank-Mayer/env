package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Frank-Mayer/env/internal"
	"github.com/atotto/clipboard"
)

func main() {
	// check for length of arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/get [profile]")
		os.Exit(1)
	}

	// get the profile name
	profile := os.Args[1]

	// get the profile keys
	profileKeysFile := "main.json"
	f, err := os.Open(profileKeysFile)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not open file %s", profileKeysFile), err))
	}
	defer f.Close()
	enc, err := io.ReadAll(f)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not read file %s", profileKeysFile), err))
	}

	key := internal.UserString("Enter the master key")
	binKey, err := internal.Atob(key)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not decode key"), err))
	}

	dec, err := internal.Decrypt(binKey, enc)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not decrypt"), err))
	}

	profileKeys := make(map[string]string)
	err = json.Unmarshal(dec, &profileKeys)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not unmarshal main.json"), err))
	}

	// get the profile key
	profileKey, ok := profileKeys[profile]
	if !ok {
		fmt.Printf("Profile %s not found\n", profile)
		os.Exit(1)
	}
	profileKeyBin, err := internal.Atob(profileKey)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not decode key"), err))
	}

	// get the profile content
	profileFile := filepath.Join("profiles", fmt.Sprintf("%s.env", profile))
	f, err = os.Open(profileFile)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not open file %s", profileFile), err))
	}
	defer f.Close()

	enc, err = io.ReadAll(f)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not read file %s", profileFile), err))
	}

	dec, err = internal.Decrypt(profileKeyBin, enc)
	if err != nil {
		panic(errors.Join(fmt.Errorf("could not decrypt"), err))
	}

	switch internal.UserChoise("What do you want to do?", "Print key", "Print content", "Copy key to clipboard", "Copy content to clipboard") {
	case 1:
		fmt.Println(profileKey)
	case 2:
		fmt.Println(string(dec))
	case 3:
		err = clipboard.WriteAll(profileKey)
		if err != nil {
			panic(errors.Join(fmt.Errorf("could not copy to clipboard"), err))
		}
	case 4:
		err = clipboard.WriteAll(string(dec))
		if err != nil {
			panic(errors.Join(fmt.Errorf("could not copy to clipboard"), err))
		}
	}
}
