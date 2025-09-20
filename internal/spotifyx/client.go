package spotifyx

import (
	"context"

	"github.com/zmb3/spotify/v2"
	auth "github.com/zmb3/spotify/v2/auth"

	"spotify-cli/internal/logx"
	"spotify-cli/internal/store"
)

// ClientFromDisk restores a spotify.Client using the saved token
func ClientFromDisk(clientID string) (*spotify.Client, error) {
	tok, err := store.LoadToken()
	if err != nil {
		return nil, logx.Err("not logged in (run: spotplay login <CLIENT_ID>): %w", err)
	}
	spAuth := auth.New(
		auth.WithClientID(clientID),
	)
	httpClient := spAuth.Client(context.Background(), tok)
	return spotify.New(httpClient), nil
}
