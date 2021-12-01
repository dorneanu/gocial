package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type IdentityRepository interface {
	Add(Identity) error
	GetByProvider(string) (Identity, error)
	Save() error
	Load() error
}

// FileIdentityRepository implements IdentityRepository interface
type FileIdentityRepository struct {
	BasePath   string
	identities []Identity
}

func NewFileIdentityRepo(path string) *FileIdentityRepository {
	return &FileIdentityRepository{
		BasePath:   path,
		identities: make([]Identity, 0),
	}
}

func (fr *FileIdentityRepository) Add(id Identity) error {
	fr.identities = append(fr.identities, id)
	return nil
}

func (fr *FileIdentityRepository) GetByProvider(provider string) (Identity, error) {
	for _, id := range fr.identities {
		if id.Provider == provider {
			return id, nil
		}
	}
	return Identity{}, fmt.Errorf("Couldn't find identity")
}

func (fr *FileIdentityRepository) Save() error {
	f, err := os.Create(fr.BasePath)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(fr.identities, "", "\t")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, bytes.NewReader(b))
	return err
}

func (fr *FileIdentityRepository) Load() error {
	f, err := os.Open(fr.BasePath)
	if err != nil {
		return fmt.Errorf("Couldn't open file: %s", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&fr.identities)
	if err != nil {
		return fmt.Errorf("Couldn't unmarshalize data: %s", err)
	}
	return nil
}
