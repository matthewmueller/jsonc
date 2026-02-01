package jsonc

import (
	"encoding/json"
	"os"

	"github.com/snorwin/jsonpatch"
)

func Unmarshal(data []byte, out any) error {
	// Don't modify the original data
	copied := make([]byte, len(data))
	copy(copied, data)

	// Standardize the JSONC data to JSON
	data, err := Standardize(copied)
	if err != nil {
		return err
	}

	// Unmarshal the standardized JSON into the output structure
	return json.Unmarshal(data, out)
}

// Patch modifies the original jsonc data with the changes, preserving comments
// and formatting.
func Patch[JSON any](prev []byte, next JSON) ([]byte, error) {
	// Parse the original JSONC file
	huValue, err := Parse(prev)
	if err != nil {
		return nil, err
	}

	// Unmarshal jsonc into json
	var old JSON
	if err := Unmarshal(prev, &old); err != nil {
		return nil, err
	}

	// Create the JSON patch between the old and new schema
	patches, err := jsonpatch.CreateJSONPatch(next, old)
	if err != nil {
		return nil, err
	}

	// Apply the patches to the JSON if there are any
	if patches.Len() > 0 {
		indentedPatch, err := json.MarshalIndent(patches.List(), "", "  ")
		if err != nil {
			return nil, err
		}
		// Patch the original JSON with the new JSON
		if err := huValue.Patch(indentedPatch); err != nil {
			return nil, err
		}
	}

	// Format the patched JSONC value
	huValue.Format()

	// Return the modified JSONC as bytes
	return []byte(huValue.String()), nil
}

// WriteFile writes the JSONC representation of data to the given path,
// preserving comments and formatting if the file already exists.
func WriteFile[JSON any](path string, data JSON) error {
	old, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			data, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return err
			}
			return os.WriteFile(path, data, 0644)
		}
		return err
	}
	patched, err := Patch(old, data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, patched, 0644)
}
