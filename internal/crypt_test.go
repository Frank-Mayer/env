package internal_test

import (
	"fmt"
	"testing"

	"github.com/Frank-Mayer/env/internal"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	tests := []string{
		"Hello, World!",
		"This is a test",
		"This is a test with a longer string",
		"FOO=BAR",
		"FOO=https://example.com/foo/bar?baz=qux",
		"FOO=https://example.com/foo/bar?baz=qux&quux=quuz\nBAR=qux\nQUX=quuz",
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			t.Parallel()
			key, err := internal.NewPassword()
			if err != nil {
				t.Error(err)
				return
			}

			plaintext := []byte(test)
			ciphertext, err := internal.Encrypt(key, plaintext)
			if err != nil {
				t.Error(err)
				return
			}
			plaintext2, err := internal.Decrypt(key, ciphertext)
			if err != nil {
				t.Error(err)
				return
			}
			if string(plaintext) != string(plaintext2) {
				t.Errorf("Plaintext and decrypted plaintext do not match: %s != %s", plaintext, plaintext2)
				return
			}
			if string(ciphertext) == string(plaintext) {
				t.Errorf("Ciphertext and plaintext are the same: %s == %s", ciphertext, plaintext)
				return
			}
		})
	}
}

func TestB64(t *testing.T) {
	t.Parallel()

	tests := []string{
		"Hello, World!",
		"This is a test",
		"This is a test with a longer string",
		"FOO=BAR",
		"FOO=https://example.com/foo/bar?baz=qux",
		"FOO=https://example.com/foo/bar?baz=qux&quux=quuz\nBAR=qux\nQUX=quuz",
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d %s", i, test), func(t *testing.T) {
			t.Parallel()
			plaintext := []byte(test)
			b64 := internal.Btoa(plaintext)
			plaintext2, err := internal.Atob(b64)
			if err != nil {
				t.Error(err)
				return
			}
			if string(plaintext) != string(plaintext2) {
				t.Errorf("Plaintext and decoded plaintext do not match: %s != %s", plaintext, plaintext2)
				return
			}
			if b64 == string(plaintext) {
				t.Errorf("Base64 and plaintext are the same: %s == %s", b64, plaintext)
				return
			}
		})
	}
}
