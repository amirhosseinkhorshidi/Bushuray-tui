## Bushuray-tui

![screenshot](./images/screenshot.png)

A keyboard-driven Xray proxy client for the terminal. Built with [bubbletea](https://github.com/charmbracelet/bubbletea) and powered by [Bushuray-core](https://github.com/Keivan-sf/Bushuray-core).

## Installation

Download the latest release from the [releases page](https://github.com/Keivan-sf/Bushuray-tui/releases), then copy both binaries:

```bash
sudo cp bushuray /usr/local/bin/bushuray
sudo cp bushuray-core /usr/local/bin/bushuray-core
```

Then run from anywhere:

```bash
bushuray
```

`bushuray-core` is looked up in the following paths (in order):

1. `./bushuray-core`
2. `./bin/bushuray-core`
3. `/usr/bin/bushuray-core`
4. `/usr/local/bin/bushuray-core`

## Keybindings

| Key | Action |
|---|---|
| `enter` | Connect / disconnect |
| `a` | Add subscription or proxy |
| `t` | Test current proxy |
| `T` | Test all proxies in subscription |
| `U` | Update subscription |
| `d` / `del` | Delete proxy |
| `D` | Delete subscription |
| `C` | Change theme |
| `↑↓` | Navigate |
| `←→` | Switch subscription tab |
| `?` | Help menu |
| `q` | Exit |

## Themes

| Name | Description |
|---|---|
| Matrix | hack the planet |
| Ocean | deep and cold |
| Nebula | lost in space |
| Ember | burn it down |
| Sakura | soft but deadly |
| Sunset | golden hour forever |

Press `C` to cycle through themes. The selected theme is saved to config automatically.

## Configuration

Config is stored at `~/.config/bushuray/config.json` and created automatically on first run.

```json
{
  "socks-port": 3090,
  "http-port": 3091,
  "core-tcp-port": 4897,
  "test-port-range": {
    "start": 3095,
    "end": 30120
  },
  "no-background": false,
  "theme": "Matrix"
}
```

| Field | Default | Description |
|---|---|---|
| `socks-port` | `3090` | Local SOCKS5 proxy port |
| `http-port` | `3091` | Local HTTP proxy port |
| `core-tcp-port` | `4897` | Port used to connect to `bushuray-core` |
| `test-port-range` | `3095–30120` | Port range used for latency testing |
| `no-background` | `false` | Disable background color (useful for transparent terminals) |
| `theme` | `Matrix` | Active theme name |

## Updating Xray

The Xray binary is located at `/opt/bushuray/bin/xray`. To update:

```bash
# Check latest version
curl -s https://api.github.com/repos/XTLS/Xray-core/releases/latest | grep tag_name

# Download and replace (update version number accordingly)
wget https://github.com/XTLS/Xray-core/releases/download/vX.X.X/Xray-linux-64.zip
unzip Xray-linux-64.zip
sudo cp xray /opt/bushuray/bin/xray
sudo chmod +x /opt/bushuray/bin/xray
```

## Debugging

```bash
tail -f debug.log       # TUI logs
tail -f core-debug.log  # bushuray-core logs
```
