<div align="center" id="top">

# Spotify-CLI

</div>
<p align="center">
  <img alt="Top language" src="https://img.shields.io/github/languages/top/yourpov/Spotify-CLI?color=56BEB8">
  <img alt="Language count" src="https://img.shields.io/github/languages/count/yourpov/Spotify-CLI?color=56BEB8">
  <img alt="Repository size" src="https://img.shields.io/github/repo-size/yourpov/Spotify-CLI?color=56BEB8">
  <img alt="License" src="https://img.shields.io/github/license/yourpov/Spotify-CLI?color=56BEB8">
</p>

---

## About

**Spotify-CLI** is a CLI written in Go for managing **Spotify playlists**

| Features                              | Description                                |
|---------------------------------------|--------------------------------------------|
| Login via Spotify OAuth               | Authorization Code flow with token storage |
| List your playlists                   | Fetch and display your playlists           |
| Create a playlist                     | Public/private with optional description   |
| Delete a playlist                     | Delete a playlist                          |
| Add / remove tracks by search         | Search by track query and modify playlists |
| Share a playlist                      | Prints an open.spotify.com share link      |


---

## Tech Stack

- [Go](https://go.dev/) (1.22+)
- [zmb3/spotify](https://github.com/zmb3/spotify) Web API client

---

## Setup

```bash
# Clone & enter project
git clone https://github.com/yourpov/Spotify-CLI
cd Spotify-CLI

# Install deps
go mod tidy

# build
go build -o Spotify-CLI.exe
```

### Create a Spotify app

- Open the Spotify Developer Dashboard
- Create an app and copy Client ID & Client Secret
- In app settings, add this Redirect URI: ```http://127.0.0.1:9876/callback```

### Configure credentials

Create a file at **`/config/config.json:`**
```json
{
  "client_id": "YOUR_CLIENT_ID",
  "client_secret": "YOUR_CLIENT_SECRET"
}
```

### Run

```go
# First-time login (stores token in your user config dir)
Spotify-CLI.exe login

# List playlists
Spotify-CLI.exe playlist list
```

## Commands

| Command                                                                             | Description                                |
|-------------------------------------------------------------------------------------|--------------------------------------------|
| `Spotify-CLI.exe login`                                                             | Authenticate with Spotify using OAuth      |
| `Spotify-CLI.exe playlist list`                                                     | List all your playlists                    |
| `Spotify-CLI.exe playlist create "<name>" [--private] [--desc "text"]`              | Create a new playlist                      |
| `Spotify-CLI.exe playlist delete "<playlist-id-or-name>"`                           | Delete a playlist                          |
| `Spotify-CLI.exe playlist add <playlist-id-or-name> "<track query>" ["more"...]`    | Add one or more tracks to a playlist       |
| `Spotify-CLI.exe playlist remove <playlist-id-or-name> "<track query>" ["more"...]` | Remove one or more tracks from a playlist  |
| `Spotify-CLI.exe playlist share <playlist-id-or-name>`                              | Print a shareable Spotify link             |

## Examples

```sh
# Login using client_id/secret from config/config.json
./Spotify-CLI.exe login

# Show all your playlists
./Spotify-CLI.exe playlist list

# Create a playlist
./Spotify-CLI.exe playlist create "late night"

# Add tracks by search (ID/URI also supported)
./Spotify-CLI.exe playlist add "late night" "playboi carti stop breathing" "central cee doja"

# Remove tracks
./Spotify-CLI.exe playlist remove "late night" "central cee doja"

# Get a shareable link
./Spotify-CLI.exe playlist share "late night"

# Delete a playlist
./Spotify-CLI.exe playlist delete "late night"
```

## Troubleshooting

- **“missing client_id or client_secret”**  
  Make sure **`config/config.json`** exists and has both fields

- **Login opens browser but errors “state mismatch”**  
  Make sure your app’s Redirect URI is `http://127.0.0.1:9876/callback`

- **403 PREMIUM_REQUIRED on playback endpoints**  
  Spotify restricts some playback APIs to Premium

- **Rate limits / network errors**  
  Re-run the command; avoid rapid loops against the API


## Showcase

![Preview](preview.gif)