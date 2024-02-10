package internal

import (
	posixpath "path"
	"path/filepath"
	"strings"
)

type Config struct {
	Schema   string     `json:"$schema,omitempty"`
	All      []Variable `json:"all" jsonschema:"title=All,description=Variables that are known to all profiles,required=true"`
	Profiles []Profile  `json:"profiles" jsonschema:"title=Profiles,description=Profiles of environment variables,required=true"`
}

func configSchema() string {
	p, _ := filepath.Abs(filepath.Join("schemas", "config.json"))
	p = "file://" + posixpath.Clean(strings.ReplaceAll(p, "\\", "/"))
	return p
}

type Profile struct {
	Name      string     `json:"name" jsonschema:"title=Name,description=Name of the profile,required=true,example=dev,example=prod"`
	Variables []Variable `json:"variables" jsonschema:"title=Variables,description=Variables of the profile,required=true"`
}

type Variable struct {
	Key   string `json:"key" jsonschema:"title=Key,description=Key of the variable,required=true,example=DB_HOST,pattern=^[A-Z_]+$"`
	Value string `json:"value" jsonschema:"title=Value,description=Value of the variable,required=true,example=localhost"`
}
