package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

const (
	appDir   = "spotplay"
	tokenFN  = "token.json"
	permDir  = 0o755
	permFile = 0o600
)

func tokenPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, appDir)
	if err := os.MkdirAll(dir, permDir); err != nil {
		return "", err
	}
	return filepath.Join(dir, tokenFN), nil
}

func SaveToken(tok *oauth2.Token) error {
	p, err := tokenPath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, permFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(tok)
}

func LoadToken() (*oauth2.Token, error) {
	p, err := tokenPath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var tok oauth2.Token
	if err := json.Unmarshal(b, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}
