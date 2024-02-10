package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Frank-Mayer/env/internal"
	"github.com/invopop/jsonschema"
)

func main() {
	config := internal.Config{}
	schema := jsonschema.Reflect(&config)
	fn := "config.json"
	err := os.MkdirAll("schemas", 0777)
	if err != nil {
		panic(err)
	}
	fp := filepath.Join("schemas", fn)
	f, err := os.Create(fp)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	if err := enc.Encode(schema); err != nil {
		panic(err)
	}
}
