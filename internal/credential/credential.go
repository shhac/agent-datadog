package credential

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shhac/agent-dd/internal/config"
)

const keychainSentinel = "__KEYCHAIN__"

type Credential struct {
	APIKey          string `json:"api_key"`
	AppKey          string `json:"app_key"`
	KeychainManaged bool   `json:"keychain_managed,omitempty"`
}

type credentialEntry struct {
	APIKey          string `json:"api_key"`
	AppKey          string `json:"app_key"`
	KeychainManaged bool   `json:"keychain_managed,omitempty"`
}

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("organization credential %q not found", e.Name)
}

func credentialsPath() string {
	return filepath.Join(config.ConfigDir(), "credentials.json")
}

func readIndex() (map[string]credentialEntry, error) {
	data, err := os.ReadFile(credentialsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]credentialEntry), nil
		}
		return nil, err
	}
	var index map[string]credentialEntry
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}
	if index == nil {
		index = make(map[string]credentialEntry)
	}
	return index, nil
}

func writeIndex(index map[string]credentialEntry) error {
	dir := config.ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credentialsPath(), append(data, '\n'), 0o600)
}

func Store(name string, cred Credential) (string, error) {
	index, err := readIndex()
	if err != nil {
		return "", err
	}

	storage := "file"
	entry := credentialEntry{
		APIKey: cred.APIKey,
		AppKey: cred.AppKey,
	}

	if err := keychainStore(name, cred.APIKey, cred.AppKey); err == nil {
		entry.APIKey = keychainSentinel
		entry.AppKey = keychainSentinel
		entry.KeychainManaged = true
		storage = "keychain"
	}

	index[name] = entry
	if err := writeIndex(index); err != nil {
		return "", err
	}
	return storage, nil
}

func Get(name string) (*Credential, error) {
	index, err := readIndex()
	if err != nil {
		return nil, err
	}
	entry, ok := index[name]
	if !ok {
		return nil, &NotFoundError{Name: name}
	}

	cred := &Credential{
		APIKey:          entry.APIKey,
		AppKey:          entry.AppKey,
		KeychainManaged: entry.KeychainManaged,
	}

	if entry.KeychainManaged {
		if apiKey, appKey, err := keychainGet(name); err == nil {
			cred.APIKey = apiKey
			cred.AppKey = appKey
		}
	}

	return cred, nil
}

func Remove(name string) error {
	index, err := readIndex()
	if err != nil {
		return err
	}
	entry, ok := index[name]
	if !ok {
		return &NotFoundError{Name: name}
	}

	if entry.KeychainManaged {
		keychainDelete(name)
	}

	delete(index, name)
	return writeIndex(index)
}

func List() ([]string, error) {
	index, err := readIndex()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(index))
	for name := range index {
		names = append(names, name)
	}
	return names, nil
}
