package jsonc_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
	"github.com/matthewmueller/jsonc"
)

const input = `// This is a schema file. It defines the server configuration for your project.
{
  // Application name
  "name": "",
  // Urls for the app (leave empty for local development)
  "urls": [
    "hello.example.com"
  ]
  // More fields
}`

const expect = `// This is a schema file. It defines the server configuration for your project.
{
  // Application name
  "name": "modified",
  // Urls for the app (leave empty for local development)
  "urls": [
    "hello.example.com",
    "new.example.com"
  ]
  // More fields
}`

type Schema struct {
	Name string   `json:"name"`
	Urls []string `json:"urls"`
}

func TestReadWrite(t *testing.T) {
	is := is.New(t)

	var schema Schema
	is.NoErr(jsonc.Unmarshal([]byte(input), &schema))

	schema.Name = "modified"
	schema.Urls = append(schema.Urls, "new.example.com")

	actual, err := jsonc.Patch([]byte(input), schema)
	is.NoErr(err)

	diff.TestString(t, strings.TrimSuffix(string(actual), "\n"), expect)
}

func TestWrite(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	schema := Schema{
		Name: "MyApp",
		Urls: []string{"app.example.com"},
	}
	expect, err := json.MarshalIndent(schema, "", "  ")
	is.NoErr(err)

	is.NoErr(jsonc.WriteFile(filepath.Join(dir, "schema.jsonc"), schema))
	actual, err := os.ReadFile(filepath.Join(dir, "schema.jsonc"))
	is.NoErr(err)

	diff.TestString(t, string(actual), string(expect))
}

const expandedInput = `// This is a schema file. It defines the server configuration for your project.
{
  // Application name
  "name": "${APP_NAME}",
  // Urls for the app (leave empty for local development)
  "urls": [
    "hello.example.com",
		"${BASE_URL}"
  ]
  // More fields
}`

const expandedExpect = `// This is a schema file. It defines the server configuration for your project.
{
  // Application name
  "name": "${APP_NAME}",
  // Urls for the app (leave empty for local development)
  "urls": [
    "hello.example.com",
    "${BASE_URL}",
    "new.example.com"
  ]
  // More fields
}`

func TestExpanded(t *testing.T) {
	is := is.New(t)

	schema, patches, err := jsonc.UnmarshalExpanded[*Schema]([]byte(expandedInput), func(key string) string {
		switch key {
		case "APP_NAME":
			return "MyApp"
		case "BASE_URL":
			return "base.example.com"
		default:
			return ""
		}
	})
	is.NoErr(err)

	schema.Urls = append(schema.Urls, "new.example.com")

	actual, err := jsonc.PatchExpanded([]byte(expandedInput), schema, patches)
	is.NoErr(err)

	diff.TestString(t, strings.TrimSuffix(string(actual), "\n"), expandedExpect)
}
