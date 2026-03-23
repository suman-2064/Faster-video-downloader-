<div align="center">

# ⚡ FASTER

### High-Speed Video & Audio Downloader for Termux

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![yt-dlp](https://img.shields.io/badge/yt--dlp-latest-FF0000?style=flat-square)](https://github.com/yt-dlp/yt-dlp)
[![Termux](https://img.shields.io/badge/Termux-compatible-1F1F1F?style=flat-square&logo=android)](https://termux.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Android%20%7C%20Linux-blue?style=flat-square)]()

```
  ███████╗ █████╗ ███████╗████████╗███████╗██████╗
  ██╔════╝██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗
  █████╗  ███████║███████╗   ██║   █████╗  ██████╔╝
  ██╔══╝  ██╔══██║╚════██║   ██║   ██╔══╝  ██╔══██╗
  ██║     ██║  ██║███████║   ██║   ███████╗██║  ██║
  ╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
```

**Download videos and audio from 1,000+ sites — directly from your Android terminal.**

[Features](#-features) · [Installation](#-installation) · [Usage](#-usage) · [Share Integration](#-share-integration) · [FAQ](#-faq)

</div>

---

## ✨ Features

| Feature | Details |
|---|---|
| 🎬 **Video Download** | 4K, 1440p, 1080p, 720p, 480p, 360p |
| 🎵 **Audio Download** | MP3 (320/128kbps), M4A, OPUS, WAV |
| 📋 **Format Listing** | Full `yt-dlp -F` output with color-coded rows |
| 🔗 **Share Integration** | Share any URL from Chrome/Instagram → Termux → faster |
| 🎨 **Beautiful TUI** | Color-coded menus, spinner, clean layout |
| ⚡ **Compiled Binary** | Built with Go — fast startup, zero runtime deps |
| 🔄 **Auto-Retry** | 5 fragment retries, 30s socket timeout |
| 📁 **Smart Save Path** | `/sdcard/Download/Faster/` (shared storage) |
| 🌐 **1,000+ Sites** | YouTube, Instagram, TikTok, Facebook, Twitter/X, Vimeo & more |

---

## 📱 Supported Platforms

- ✅ **Termux** on Android 7.0+ (primary target)
- ✅ Linux (Ubuntu, Debian, Arch)
- ✅ macOS (with Homebrew Go + yt-dlp)

---

## 🚀 Installation

### Prerequisites — one-time setup

Open Termux and run these commands **in order**:

#### 1. Update Termux packages

```bash
pkg update && pkg upgrade -y
```

#### 2. Install system dependencies

```bash
pkg install golang ffmpeg python git curl -y
```

> **Why each package?**
> - `golang` — compiles the `faster` binary
> - `ffmpeg` — merges separate video+audio streams (required for 1080p+)
> - `python` — runs `yt-dlp`
> - `git` — clones this repo
> - `curl` — optional, useful for testing connectivity

#### 3. Install yt-dlp

```bash
pip install yt-dlp
```

Verify it works:

```bash
yt-dlp --version
```

---

### Clone & Build

#### 4. Clone the repository

```bash
git clone https://github.com/YOUR_USERNAME/faster.git
cd faster
```

#### 5. Run the setup script

```bash
bash setup.sh
```

The script automatically:

- Compiles `main.go` with `go build`
- Copies the binary to `$PREFIX/bin/faster` (globally accessible)
- Installs `termux-url-opener` for Share integration
- Adds a shell alias

#### 6. Reload your shell

```bash
source ~/.bashrc
# or, if you use zsh:
source ~/.zshrc
```

#### 7. Verify installation

```bash
faster --check
```

---

### Manual Build (optional)

If you prefer to build manually:

```bash
cd faster

# Initialise Go module
go mod init faster

# Build optimised binary
go build -ldflags="-s -w" -o faster main.go

# Install globally
cp faster $PREFIX/bin/faster
chmod +x $PREFIX/bin/faster
```

---

## 🎯 Usage

### Interactive Mode

Just type:

```bash
faster
```

You'll see:

```
  ➜  Paste your YouTube / Instagram / Facebook / TikTok URL
  ➜  Type q to quit

  ❯ URL: https://www.youtube.com/watch?v=dQw4w9WgXcQ

  ⠏ Fetching available formats…
  ✓ Fetching available formats…

  ──────────────────────────────────────────────────────
  📋 Available Formats
  ──────────────────────────────────────────────────────
  ID   EXT   RESOLUTION  FPS   FILESIZE   TBR  ...
  ...
  ──────────────────────────────────────────────────────

  ┌─────────────────────────────────────┐
  │       What do you want to save?     │
  ├─────────────────────────────────────┤
  │  [1] 🎬  Video  (with audio)        │
  │  [2] 🎵  Audio  (music/podcast)     │
  │  [3] 🎞   Custom (enter format ID)  │
  │  [q] ✕   Quit                       │
  └─────────────────────────────────────┘
```

### Direct URL Mode

Pass the URL as an argument to skip the prompt:

```bash
faster "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

### All Commands

| Command | Description |
|---|---|
| `faster` | Start interactive mode |
| `faster <URL>` | Download with URL pre-filled |
| `faster --check` | Verify all dependencies |
| `faster --version` | Show version number |
| `faster --help` | Show help |

---

## 🔗 Share Integration

**faster** installs a `termux-url-opener` script that lets you share any video link directly into Termux.

### How to use it

1. Open YouTube, Instagram, TikTok, or any supported site in your browser
2. Tap **Share**
3. Select **Termux** from the share sheet
4. `faster` opens automatically with the URL pre-filled
5. Choose your format and quality → download starts

### Manual setup (if setup.sh was not used)

```bash
# Copy the url-opener to Termux's bin
cp termux-url-opener $PREFIX/bin/termux-url-opener
chmod +x $PREFIX/bin/termux-url-opener
```

---

## 📂 Where are files saved?

| Platform | Path |
|---|---|
| Android (Termux) | `/sdcard/Download/Faster/` |
| Linux / macOS | `~/Downloads/Faster/` |

Files are named using the video title automatically.

---

## 🎬 Video Quality Guide

| Option | Format | Notes |
|---|---|---|
| 4K / 2160p | `bestvideo[height<=2160]+bestaudio` | Requires ffmpeg |
| 1440p | `bestvideo[height<=1440]+bestaudio` | Requires ffmpeg |
| 1080p Full HD | `bestvideo[height<=1080]+bestaudio` | Most popular |
| 720p HD | `bestvideo[height<=720]+bestaudio` | Good balance |
| 480p | `bestvideo[height<=480]+bestaudio` | Data-saving |
| 360p | `bestvideo[height<=360]+bestaudio` | Very small |
| Best auto | `bestvideo+bestaudio/best` | Let yt-dlp decide |
| Custom | User-entered format ID | From `yt-dlp -F` output |

---

## 🎵 Audio Format Guide

| Option | Format | Bitrate | Notes |
|---|---|---|---|
| MP3 High | MP3 | 320kbps | Best quality MP3 |
| MP3 Standard | MP3 | 128kbps | Smaller file |
| M4A | M4A | Best | iTunes/Apple compatible |
| OPUS | OPUS | Best | Smallest, great quality |
| WAV | WAV | Lossless | Largest file |
| Auto | Best | — | yt-dlp chooses |

---

## 🔧 Updating

### Update yt-dlp (do this regularly)

```bash
pip install --upgrade yt-dlp
```

### Update faster

```bash
cd faster
git pull origin main
bash setup.sh
```

---

## 🏗 Project Structure

```
faster/
├── main.go             # Core Go application
├── setup.sh            # One-command installer
├── go.mod              # Go module (auto-generated)
├── termux-url-opener   # Share integration script
└── README.md           # This file
```

---

## ❓ FAQ

**Q: I get `yt-dlp: command not found`**
```bash
pip install yt-dlp
# If that fails:
pip3 install yt-dlp
```

**Q: Downloads fail with "HTTP Error 429"**
YouTube is rate-limiting you. Wait a few minutes and try again. You can also use cookies:
```bash
yt-dlp --cookies-from-browser chrome <URL>
```

**Q: 1080p video has no audio**
Make sure `ffmpeg` is installed — it's needed to merge video+audio streams:
```bash
pkg install ffmpeg -y
```

**Q: The binary says "exec format error"**
You need to recompile for your architecture. Run `setup.sh` again from inside Termux.

**Q: Shared storage not available (`/sdcard` missing)**
Grant Termux storage permission:
```bash
termux-setup-storage
```

**Q: How do I use cookies for private/age-restricted videos?**
```bash
yt-dlp --cookies cookies.txt <URL>
```
Export cookies from your browser using a cookies.txt extension.

**Q: Can I use faster on a rooted device?**
Yes — it works the same way. No root is required or used.

---

## 🌐 Supported Sites (partial list)

- YouTube & YouTube Music
- Instagram (Reels, Stories, Posts)
- Facebook (Videos, Reels)
- TikTok
- Twitter / X
- Vimeo
- Dailymotion
- Reddit
- Twitch (VODs & clips)
- SoundCloud
- Bandcamp
- And **1,000+ more** via [yt-dlp's supported sites](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md)

---

## 🛡 Legal Notice

This tool uses [yt-dlp](https://github.com/yt-dlp/yt-dlp) under the hood. Only download content you have the right to download. Respect copyright laws in your country. This project is intended for personal use with content you own or that is freely licensed.

---

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">

Made with ❤️ for the Termux community

**[⭐ Star this repo](https://github.com/YOUR_USERNAME/faster)** if faster saves you time!

</div>
