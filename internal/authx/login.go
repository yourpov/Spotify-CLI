package authx

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	logger "github.com/yourpov/logrite"
	"github.com/zmb3/spotify/v2"
	auth "github.com/zmb3/spotify/v2/auth"

	"spotify-cli/internal/logx"
	"spotify-cli/internal/store"
)

const redirectURL = "http://127.0.0.1:9876/callback"

func b64url(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}
func randBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

// Login performs Authorization Code OAuth, saves token, and returns a spotify.Client
func Login(clientID, clientSecret string) (*spotify.Client, error) {
	if clientSecret == "" {
		return nil, logx.Err("client secret required")
	}

	ctx := context.Background()
	state := b64url(randBytes(12))

	spAuth := auth.New(auth.WithClientID(clientID), auth.WithClientSecret(clientSecret), auth.WithRedirectURL(redirectURL), auth.WithScopes("playlist-read-private", "playlist-modify-private", "playlist-modify-public", "user-read-private"))

	authURL := spAuth.AuthURL(state)

	reqCh := make(chan *http.Request, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:              "127.0.0.1:9876",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			errCh <- logx.Err("state mismatch")
			http.Error(w, "state mismatch", http.StatusBadRequest)
			return
		}
		if e := r.URL.Query().Get("error"); e != "" {
			errCh <- logx.Err("spotify error: %s", e)
			http.Error(w, e, http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("code") == "" {
			errCh <- logx.Err("missing code")
			http.Error(w, "missing code", http.StatusBadRequest)
			return
		}
		logger.Success("Login complete. You can close this tab.")
		reqCh <- r
	})

	logger.Custom("🌐", "open this url to authorize:\n", authURL, logger.Black, logger.BgWhite)
	go srv.ListenAndServe()

	var req *http.Request
	select {
	case req = <-reqCh:
	case err := <-errCh:
		return nil, err
	case <-time.After(2 * time.Minute):
		return nil, logx.Err("timed out waiting for OAuth callback")
	}
	_ = srv.Shutdown(ctx)

	tok, err := spAuth.Token(ctx, state, req)
	if err != nil {
		return nil, err
	}
	if err := store.SaveToken(tok); err != nil {
		logger.Warn("warning: failed to save token:", err)
	}

	httpClient := spAuth.Client(ctx, tok)
	logger.Success("Logged in. Token saved.")
	return spotify.New(httpClient), nil
}
