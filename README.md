# Goify

## This project ⭐

* is written entirely in Golang
* doesn't utilize any external APIs
* uses Playwright script to scrape tracks from provided Spotify URL and downloads them from a convertor page
* downloads the tracks in mp3 format

## Quickstart 🚀

To run goify and download the desired spotify playlist, album or track:
* on Spotify copy link of any track, album or playlist (as long as it is public)
* run the command goify with -download and the desired Spotify URL (e.g. https://open.spotify.com/playlist/...)
```shell
goify -download "https://open.spotify.com/playlist/..."
```

To maximize download speed for your PC you can also tweak other flags:
```shell
goify -download "https://open.spotify.com/playlist/..." -workers=3 -headless=0
```

Use -help to list all the available flags and their purpose
```shell
goify -help
```

### All available flags

| Option        | Type   | Description |
|--------------|--------|-------------|
| `-atimeout`  | `int`  | Time in ms to wait for an action performed on a locator (default: `3000`) |
| `-convertor` | `string` | URL of a Spotify-to-audio converter without captcha (default: `"https://spotidown.app/"`) |
| `-dir`       | `string` | Directory on this PC to download the track(s) to (default: `"./"`) |
| `-download`  | `string` | Spotify URL of a track, album, or playlist |
| `-dtimeout`  | `int`  | Time in ms to wait for a download before retrying (default: `5000`) |
| `-headless`  | `bool` | Hide browser windows or not (default: `true`) |
| `-sleep`     | `int`  | Time in ms for the program to sleep between repeating actions (default: `1000`) |
| `-uagent`    | `string` | Provide a custom user agent (default: `"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"`) |
| `-workers`   | `int`  | Maximum number of browser contexts attempting downloads simultaneously (default: `5`) |

## How to install goify 🔥

* download and install Go for your OS from their official website: https://go.dev/dl/ 
* use GIT to pull this repo to your PC or download it from here in a .zip file

### On Windows
* run the batch file which compiles the app and moves it to a designated directory
* sometimes the script refuses to add the directory to PATH, in which case please do so manually
```bash
.\install_goify.bat
```
### On Unix

* run the shell file which functions the same as the one for windows
```shell
.\install_goify.sh
```

## Disclaimer

* this tool is intended for experimental and educational purposes only 
* I do not endorse or encourage downloading copyrighted content from Spotify or any other streaming service
* users at your own discretion