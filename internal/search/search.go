package search

import (
	"context"

	"github.com/zmb3/spotify/v2"

	"spotify-cli/internal/logx"
)

// TrackID returns the first track for a search
func TrackID(ctx context.Context, cli *spotify.Client, q string) (spotify.ID, error) {
	res, err := cli.Search(ctx, q, spotify.SearchTypeTrack, spotify.Limit(1))
	if err != nil {
		return "", err
	}
	if res.Tracks == nil || len(res.Tracks.Tracks) == 0 {
		return "", logx.Err("no track for %q", q)
	}
	return res.Tracks.Tracks[0].ID, nil
}
