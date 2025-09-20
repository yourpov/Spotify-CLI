package playlist

import (
	"context"
	"strings"

	logger "github.com/yourpov/logrite"
	"github.com/zmb3/spotify/v2"

	"spotify-cli/internal/logx"
	"spotify-cli/internal/search"
	"spotify-cli/internal/util"
)

// List lists playlists
func List(cli *spotify.Client) error {
	ctx := context.Background()
	logger.Info("%-24s  %-34s  %6s  %-s", "ID", "NAME", "TRACKS", "OWNER")
	var offset int
	for {
		page, err := cli.CurrentUsersPlaylists(ctx, spotify.Limit(50), spotify.Offset(offset))
		if err != nil {
			return err
		}
		for _, p := range page.Playlists {
			owner := p.Owner.DisplayName
			if owner == "" {
				owner = p.Owner.ID
			}
			logger.Custom("📁", "list", "%-24s  %-34.34s  %6d  %-s", logger.Black, logger.BgWhite, string(p.ID), p.Name, p.Tracks.Total, owner)
		}

		if len(page.Playlists) < 50 {
			break
		}
		offset += 50
	}
	return nil
}

// Create makes a playlist
func Create(cli *spotify.Client, name string, priv bool, desc string) error {
	ctx := context.Background()
	me, err := cli.CurrentUser(ctx)
	if err != nil {
		return err
	}
	pl, err := cli.CreatePlaylistForUser(ctx, me.ID, name, desc, !priv, false)
	if err != nil {
		return err
	}
	vis := "public"
	if !pl.IsPublic {
		vis = "private"
	}
	logger.Custom("✅", "create", "created: %s  %s  (%s)\nshare:   https://open.spotify.com/playlist/%s", logger.Black, logger.BgGreen, pl.ID, pl.Name, vis, pl.ID)
	return nil
}

// Delete removes a playlist
func Delete(cli *spotify.Client, ident string) error {
	ctx := context.Background()
	id, name, err := resolveByIdent(ctx, cli, ident)
	if err != nil {
		return err
	}
	if err := cli.UnfollowPlaylist(ctx, id); err != nil {
		return err
	}
	logger.Custom("🗑️", "delete", "deleted playlist %s (%s)", logger.Black, logger.BgRed, id, name)
	return nil
}

// resolveByIdent identifies playlist
func resolveByIdent(ctx context.Context, cli *spotify.Client, ident string) (spotify.ID, string, error) {
	if util.LooksLikePlaylistID(ident) {
		id := util.NormalizePlaylistID(ident)
		pl, err := cli.GetPlaylist(ctx, id)
		if err == nil && pl != nil {
			return id, pl.Name, nil
		}
	}
	var offset int
	for {
		page, err := cli.CurrentUsersPlaylists(ctx, spotify.Limit(50), spotify.Offset(offset))
		if err != nil {
			return "", "", err
		}
		for _, p := range page.Playlists {
			if strings.EqualFold(p.Name, ident) {
				return p.ID, p.Name, nil
			}
		}
		if len(page.Playlists) < 50 {
			break
		}
		offset += 50
	}
	return "", "", logx.Err("playlist %q not found", ident)
}

// Add adds tracks to Spotify
func Add(cli *spotify.Client, ident string, queries []string) error {
	ctx := context.Background()
	id, name, err := resolveByIdent(ctx, cli, ident)
	if err != nil {
		return err
	}
	var ids []spotify.ID
	for _, q := range queries {
		tid, err := search.TrackID(ctx, cli, q)
		if err != nil {
			return err
		}
		ids = append(ids, tid)
	}
	_, err = cli.AddTracksToPlaylist(ctx, id, ids...)
	if err != nil {
		return err
	}
	logger.Custom("✅", "add", "added %d track(s) to %s\n", logger.Black, logger.BgGreen, len(ids), name)
	return nil
}

// Remove removes tracks from Spotify
func Remove(cli *spotify.Client, ident string, queries []string) error {
	ctx := context.Background()
	id, name, err := resolveByIdent(ctx, cli, ident)
	if err != nil {
		return err
	}
	var ids []spotify.ID
	for _, q := range queries {
		tid, err := search.TrackID(ctx, cli, q)
		if err != nil {
			return err
		}
		ids = append(ids, tid)
	}
	_, err = cli.RemoveTracksFromPlaylist(ctx, id, ids...)
	if err != nil {
		return err
	}
	logger.Custom("🗑️", "remove", "removed %d track(s) from %s", logger.Black, logger.BgGreen, len(ids), name)
	return nil
}

// Share shares a Spotify playlist by printing out the URL
func Share(cli *spotify.Client, ident string) error {
	ctx := context.Background()
	id, _, err := resolveByIdent(ctx, cli, ident)
	if err != nil {
		return err
	}
	logger.Custom("✅", "share", "https://open.spotify.com/playlist/"+string(id), logger.BgGreen, logger.Black)
	return nil
}
