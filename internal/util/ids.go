package util

import (
	"net/url"
	"strings"

	"github.com/zmb3/spotify/v2"
)

func LooksLikePlaylistID(s string) bool {
	if strings.HasPrefix(s, "spotify:playlist:") {
		return true
	}
	if strings.HasPrefix(s, "https://open.spotify.com/playlist/") {
		return true
	}
	return len(s) >= 22 && len(s) <= 64
}

func NormalizePlaylistID(s string) spotify.ID {
	if strings.HasPrefix(s, "spotify:playlist:") {
		return spotify.ID(strings.TrimPrefix(s, "spotify:playlist:"))
	}
	if strings.HasPrefix(s, "https://open.spotify.com/playlist/") {
		u, _ := url.Parse(s)
		if u != nil {
			parts := strings.Split(strings.TrimSuffix(u.Path, "/"), "/")
			return spotify.ID(parts[len(parts)-1])
		}
	}
	return spotify.ID(s)
}
