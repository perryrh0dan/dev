package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Identity holds a Git user identity.
type Identity struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	KeyID  string `json:"keyid"`
	Active bool   `json:"active"`
}

// SaveIdentity upserts an identity and marks it as the active one.
func SaveIdentity(id Identity) error {
	path, err := identitiesPath()
	if err != nil {
		return err
	}
	return saveIdentityToPath(path, id)
}

// ActivateIdentity marks the identity with the given email as active.
func ActivateIdentity(email string) error {
	path, err := identitiesPath()
	if err != nil {
		return err
	}
	return activateIdentityInPath(path, email)
}

// LoadActiveIdentity returns the currently active identity, or nil if none is set.
func LoadActiveIdentity() (*Identity, error) {
	path, err := identitiesPath()
	if err != nil {
		return nil, err
	}
	return loadActiveIdentityFromPath(path)
}

func identitiesPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".config", "dev", "identities.json"), nil
}

func loadIdentitiesFromPath(path string) ([]Identity, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Identity{}, nil
	}
	if err != nil {
		return nil, err
	}
	var ids []Identity
	if err := json.Unmarshal(data, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func saveIdentityToPath(path string, id Identity) error {
	ids, err := loadIdentitiesFromPath(path)
	if err != nil {
		return err
	}
	// Mark all others inactive, new identity is always active
	id.Active = true
	updated := false
	for i, existing := range ids {
		ids[i].Active = false
		if existing.Email == id.Email {
			ids[i] = id
			updated = true
		}
	}
	if !updated {
		ids = append(ids, id)
	}
	return writeJSON(path, ids)
}

func activateIdentityInPath(path, email string) error {
	ids, err := loadIdentitiesFromPath(path)
	if err != nil {
		return err
	}
	found := false
	for i, id := range ids {
		ids[i].Active = id.Email == email
		if id.Email == email {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("identity %q not found", email)
	}
	return writeJSON(path, ids)
}

func loadActiveIdentityFromPath(path string) (*Identity, error) {
	ids, err := loadIdentitiesFromPath(path)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		if id.Active {
			id := id
			return &id, nil
		}
	}
	return nil, nil
}
