# i3-workspace-renamer

Simple i3 IPC script that renames Workspaces on-the-fly based on its contents.
Currently it renames based on what I use on my setup (`event_new.go`).

## How

The application subscribes to events and proceeds to rename the needed workspaces when each event arrives.
When the connection drops for whatever reason, the application will try to open a new connection every 5 seconds.
This is interesting for when you need to restart your i3wm instance and you run this application under your configuration file.

## Configuration

All the configuration happens inside the `i3wr_config.json` file.
The configuration file is reloaded when the socket is closed (for example when reloading i3wm).
Below is a sample configuration file:

```json
{
	"separator": " | ",
	"class": {
		"gimp":          "✎ Gimp",
		"clockify":      "🗒  Clockify",
		"google-chrome": "◎ Chrome",
		"st":            "▱ Terminal",
		"discord":       "🗪 Discord",
		"spotify":       "🎵 Spotify"
	},
	"window": {
		"vim":  "▤ Vim",
		"gimp": "✎ Gimp"
	}
}
```

Specific matches for applications should usually go inside the `class` option.
Scripts that run on the terminal can be also matches with the `window` option (for example `htop`).
