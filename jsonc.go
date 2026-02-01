package jsonc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/snorwin/jsonpatch"
)

// Patches represents a list of JSON patches.
type Patches = jsonpatch.JSONPatchList

func Unmarshal(data []byte, out any) error {
	// Don't modify the original data
	copied := make([]byte, len(data))
	copy(copied, data)

	// Standardize the JSONC data to JSON
	data, err := Standardize(copied)
	if err != nil {
		return err
	}

	// Create a decoder and disallow unknown fields
	decoder := json.NewDecoder(bytes.NewReader(data))
	// Disallow unknown fields
	decoder.DisallowUnknownFields()

	// Unmarshal the JSON into the schema
	if err := decoder.Decode(&out); err != nil {
		return fmt.Errorf("jsonc: failed to unmarshal data: %w", err)
	}

	return nil
}

// UnmarshalExpanded unmarshals JSONC data into the given structure, expanding any
// environment variables using the provided expander function. It also returns the
// JSON patches that were applied during the expansion.
func UnmarshalExpanded[JSON any](data []byte, expander func(key string) string) (out JSON, reverts Patches, err error) {
	if err := Unmarshal(data, &out); err != nil {
		return out, reverts, err
	}
	expandedData := os.Expand(string(data), expander)
	var expanded JSON
	if err := Unmarshal([]byte(expandedData), &expanded); err != nil {
		return out, reverts, err
	}
	reverts, err = jsonpatch.CreateJSONPatch(out, expanded)
	if err != nil {
		return out, reverts, err
	}
	return expanded, reverts, nil
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

// PatchExpanded modifies the original expanded jsonc data with the changes,
// reverting any expansions made previously, and preserving comments and
// formatting.
// TODO: right now the reverts always win even if the next value was modified
// after the original expansion. See if we can improve this behavior.
func PatchExpanded[JSON any](prev []byte, next JSON, reverts Patches) ([]byte, error) {
	// Apply the original patches to bring the schema back to the unexpanded state
	if reverts.Len() > 0 {
		nextData, err := json.Marshal(next)
		if err != nil {
			return nil, err
		}
		nextValue, err := Parse(nextData)
		if err != nil {
			return nil, err
		}
		indentedPatch, err := json.MarshalIndent(reverts.List(), "", "  ")
		if err != nil {
			return nil, err
		}
		// Patch the next JSON with the new JSON
		if err := nextValue.Patch(indentedPatch); err != nil {
			return nil, err
		}
		// Unmarshal the patched JSON back into the next
		if err := json.Unmarshal([]byte(nextValue.String()), &next); err != nil {
			return nil, err
		}
	}
	return Patch(prev, next)
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
