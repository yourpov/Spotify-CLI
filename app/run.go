package app

import (
	"flag"

	logger "github.com/yourpov/logrite"

	"spotify-cli/config"
	"spotify-cli/internal/authx"
	"spotify-cli/internal/logx"
	"spotify-cli/internal/playlist"
	"spotify-cli/internal/spotifyx"
)

func usage() {
	logger.Success(`Spotify playlist CLI

Usage:
  Spotify-CLI.exe login
  Spotify-CLI.exe playlist <subcommand>

Playlist subcommands:
  playlist list
  playlist create "<name>" [--private] [--desc "text"]
  playlist add <playlist-id-or-name> "<track query>" ["more"...]
  playlist remove <playlist-id-or-name> "<track query>" ["more"...]
  playlist share <playlist-id-or-name>
`)
}

func Run(args []string) error {
	if len(args) == 0 {
		usage()
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		return logx.Err("failed to load config: %s", err)
	}
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return logx.Err("[Config]: client_id or client_secret is empty")
	}

	switch args[0] {
	case "login":
		_, err := authx.Login(cfg.ClientID, cfg.ClientSecret)
		return err

	case "playlist":
		cli, err := spotifyx.ClientFromDisk(cfg.ClientID)
		if err != nil {
			return err
		}
		if len(args) < 2 {
			usage()
			return nil
		}
		switch args[1] {
		case "list":
			return playlist.List(cli)

		case "create":
			fs := flag.NewFlagSet("create", flag.ContinueOnError)
			priv := fs.Bool("private", false, "create as private")
			desc := fs.String("desc", "", "playlist description")
			if err := fs.Parse(args[2:]); err != nil {
				return err
			}
			left := fs.Args()
			if len(left) < 1 {
				return logx.Err(`usage: Spotify-CLI.exe playlist create "<name>" [--private] [--desc "text"]`)
			}
			return playlist.Create(cli, left[0], *priv, *desc)

		case "delete":
			if len(args) < 3 {
				return logx.Err(`usage: spotplay playlist delete <playlist-id-or-name>`)
			}
			return playlist.Delete(cli, args[2])

		case "add":
			if len(args) < 4 {
				return logx.Err(`usage: Spotify-CLI.exe playlist add <playlist-id-or-name> "<track query>" ["more"...]`)
			}
			return playlist.Add(cli, args[2], args[3:])

		case "remove":
			if len(args) < 4 {
				return logx.Err(`usage: Spotify-CLI.exe playlist remove <playlist-id-or-name> "<track query>" ["more"...]`)
			}
			return playlist.Remove(cli, args[2], args[3:])

		case "share":
			if len(args) < 3 {
				return logx.Err(`usage: Spotify-CLI.exe playlist share <playlist-id-or-name>`)
			}
			return playlist.Share(cli, args[2])

		default:
			usage()
			return nil
		}

	default:
		usage()
		return nil
	}
}
