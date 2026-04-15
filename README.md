# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a baseline scan:

```bash
portwatch start
```

Watch specific ports and get alerted when unexpected ports open or close:

```bash
portwatch start --interval 30s --alert-cmd "notify-send 'Port change detected'"
```

Run a one-time snapshot and compare against a saved baseline:

```bash
portwatch scan --save baseline.json
portwatch scan --diff baseline.json
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `60s` | How often to poll open ports |
| `--alert-cmd` | `` | Shell command to run on change |
| `--log` | `stdout` | Log output destination |
| `--ports` | all | Comma-separated ports to watch |

## How It Works

`portwatch` periodically reads the system's open port list, compares it against a known-good baseline, and triggers configurable alerts when unexpected changes are detected. It runs as a lightweight background process with minimal resource usage.

## License

MIT © 2024 yourusername