package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// writeJSON atomically writes v as JSON to path (write unique temp + rename).
// A deferred cleanup removes the temp file if any step fails before the rename.
func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	success := false
	defer func() {
		if !success {
			os.Remove(tmp)
		}
	}()

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		return err
	}
	success = true
	return nil
}
