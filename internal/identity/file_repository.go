package identity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dorneanu/gomation/internal/entity"
	"github.com/labstack/echo/v4"
)

// FileIdentityRepository implements IdentityRepository interface
type FileIdentityRepository struct {
	BasePath   string
	identities []entity.IdentityProvider
}

func NewFileIdentityRepo(path string) *FileIdentityRepository {
	return &FileIdentityRepository{
		BasePath:   path,
		identities: make([]entity.IdentityProvider, 0),
	}
}

func (fr *FileIdentityRepository) Add(id entity.IdentityProvider, c echo.Context) error {
	fr.identities = append(fr.identities, id)
	return nil
}

func (fr *FileIdentityRepository) GetByProvider(c echo.Context, provider string) (entity.IdentityProvider, error) {
	for _, id := range fr.identities {
		if id.Provider == provider {
			return id, nil
		}
	}
	return entity.IdentityProvider{}, fmt.Errorf("Couldn't find identity")
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
