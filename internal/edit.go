package internal

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
)

func EditStruct(s interface{}) (error, bool) {
	tmgDir := os.TempDir()
	r := fmt.Sprintf("%x", rand.Int())
	fn := filepath.Join(tmgDir, r+".json")

	var initContent [20]byte

	// create temporary file
	if err := func() error {
		f, err := os.Create(fn)
		if err != nil {
			return errors.Join(errors.New("failed to create temporary file"), err)
		}
		defer f.Close()
		jsonEnc, err := json.MarshalIndent(s, "", "    ")
		if err != nil {
			return errors.Join(errors.New("failed to encode struct"), err)
		}
		if _, err := f.Write(jsonEnc); err != nil {
			return errors.Join(errors.New("failed to write to temporary file"), err)
		}
		initContent = sha1.Sum(jsonEnc)
		return nil
	}(); err != nil {
		return err, false
	}

	// delete temporary file
	defer os.Remove(fn)

	// edit file by user
	if err := EditFile(fn); err != nil {
		return errors.Join(errors.New("failed to edit file"), err), false
	}

	// read file back into struct

	f, err := os.Open(fn)
	if err != nil {
		return errors.Join(errors.New("failed to open temporary file"), err), false
	}
	defer f.Close()
	// read file into string
	jsonEnc, err := os.ReadFile(fn)
	if err != nil {
		return errors.Join(errors.New("failed to read temporary file"), err), false
	}
	// check if file was changed
	newContent := sha1.Sum(jsonEnc)
	change := false
	for i := range initContent {
		if initContent[i] != newContent[i] {
			change = true
			break
		}
	}
	if !change {
		return errors.New("Content did not change"), false
	}
	// unmarshal json into struct
	if err := json.Unmarshal(jsonEnc, s); err != nil {
		return nil, false
	}
	return nil, true
}

func EditFile(file string) error {
	editor, err := editor()
	if err != nil {
		return errors.Join(errors.New("failed to get editor"), err)
	}
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Join(errors.New("failed to run editor"), err)
	}
	return nil
}

func editor() (string, error) {
	// look for specific environment variables
	if visual, ok := os.LookupEnv("VISUAL"); ok {
		return visual, nil
	}
	if editor, ok := os.LookupEnv("EDITOR"); ok {
		return editor, nil
	}
	if gitEditor, ok := os.LookupEnv("GIT_EDITOR"); ok {
		return gitEditor, nil
	}

	// look if any known editor is available
	if subl, err := exec.LookPath("subl"); err == nil && subl != "" {
		return subl, nil
	}
	if code, err := exec.LookPath("code"); err == nil && code != "" {
		return code, nil
	}
	if codium, err := exec.LookPath("codium"); err == nil && codium != "" {
		return codium, nil
	}
	if nvim, err := exec.LookPath("nvim"); err == nil && nvim != "" {
		return nvim, nil
	}

	return "", errors.New("no editor found")
}
