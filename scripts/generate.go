package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"brew-manager/pkg/types"

	"github.com/invopop/jsonschema"
)

func main() {
	outputPath := "./packages.schema.json"
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}

	ref := &jsonschema.Reflector{
		FieldNameTag:               "yaml",
		RequiredFromJSONSchemaTags: true,
	}

	schema := ref.Reflect(&types.PackageGrouped{})

	schemaStr, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err)
	}

	_, err = file.Write(schemaStr)
	if err != nil {
		panic(err)
	}
}
