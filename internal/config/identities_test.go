package config

import (
	"path/filepath"
	"testing"
)

func TestSaveAndLoadIdentity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "identities.json")

	id := Identity{Email: "me@example.com", Name: "Me", KeyID: "ABC123", Active: true}
	if err := saveIdentityToPath(path, id); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	ids, err := loadIdentitiesFromPath(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(ids) != 1 || ids[0].Email != "me@example.com" {
		t.Errorf("unexpected: %+v", ids)
	}
}

func TestSaveIdentity_OnlyOneActive(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "identities.json")

	saveIdentityToPath(path, Identity{Email: "a@example.com", Name: "A", KeyID: "K1", Active: true})
	saveIdentityToPath(path, Identity{Email: "b@example.com", Name: "B", KeyID: "K2", Active: true})

	ids, _ := loadIdentitiesFromPath(path)
	activeCount := 0
	for _, id := range ids {
		if id.Active {
			activeCount++
		}
	}
	if activeCount != 1 {
		t.Errorf("expected exactly 1 active identity, got %d", activeCount)
	}
}

func TestLoadActiveIdentity_ReturnsActive(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "identities.json")

	saveIdentityToPath(path, Identity{Email: "a@example.com", Name: "A", KeyID: "K1", Active: false})
	saveIdentityToPath(path, Identity{Email: "b@example.com", Name: "B", KeyID: "K2", Active: true})

	active, err := loadActiveIdentityFromPath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if active == nil || active.Email != "b@example.com" {
		t.Errorf("expected b@example.com, got %v", active)
	}
}

func TestActivateIdentity_SwitchesActive(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "identities.json")

	saveIdentityToPath(path, Identity{Email: "a@example.com", Name: "A", KeyID: "K1", Active: true})
	saveIdentityToPath(path, Identity{Email: "b@example.com", Name: "B", KeyID: "K2", Active: false})

	if err := activateIdentityInPath(path, "b@example.com"); err != nil {
		t.Fatalf("activate failed: %v", err)
	}

	active, _ := loadActiveIdentityFromPath(path)
	if active == nil || active.Email != "b@example.com" {
		t.Errorf("expected b@example.com active, got %v", active)
	}
}
