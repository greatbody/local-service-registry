# Local Service Registry

[中文文档](README_zh-CN.md)

A lightweight local service registry with health checks and a built-in web UI. It gives you a single place to keep track of all the web services running on your machine, so you never forget what's running and where.

## Features

- **Service registration** via REST API -- services can register themselves on startup
- **Periodic health checks** -- every 5 minutes (configurable), GET each registered URL and check for 2xx status
- **Immediate check on registration** -- newly registered services are probed right away
- **Built-in web UI** -- view all services, their health status, and manage registrations from the browser
- **SQLite persistence** -- registrations survive restarts
- **AI prompt helper** -- one-click button to copy a prompt that instructs an AI to add self-registration to any project
- **Zero configuration** -- just run the binary

## Installation

### From source

```bash
go install github.com/greatbody/local-service-registry@latest
```

### Build locally

```bash
git clone https://github.com/greatbody/local-service-registry.git
cd local-service-registry
go build -o local-service-registry .
```

## Usage

```bash
# Start with defaults (port 8500, 5-minute health check interval)
./local-service-registry

# Custom port and interval
./local-service-registry -addr :1234 -interval 30s -db /path/to/registry.db
```

Then open `http://localhost:1234` in your browser to see the web UI.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-addr` | `:8500` | HTTP listen address |
| `-db` | `registry.db` | SQLite database file path |
| `-interval` | `5m` | Health check interval |

## API

### Register a service

```bash
curl -X POST http://localhost:1234/services \
  -H 'Content-Type: application/json' \
  -d '{"name": "my-app", "url": "http://localhost:3000", "description": "My web app"}'
```

### List all services

```bash
curl http://localhost:1234/services
```

### Get a single service

```bash
curl http://localhost:1234/services/<id>
```

### Remove a service

```bash
curl -X DELETE http://localhost:1234/services/<id>
```

## Run as a macOS service

Create a LaunchAgent to start the registry on login:

```bash
cat > ~/Library/LaunchAgents/com.local-service-registry.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.local-service-registry</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/local-service-registry</string>
        <string>-addr</string>
        <string>:1234</string>
        <string>-db</string>
        <string>/Users/YOU/Library/Application Support/local-service-registry/registry.db</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/YOU/Library/Logs/local-service-registry.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/YOU/Library/Logs/local-service-registry.log</string>
</dict>
</plist>
EOF

# Load and start
launchctl load ~/Library/LaunchAgents/com.local-service-registry.plist
```

Replace `/Users/YOU` with your actual home directory.

## Self-registration from other services

The web UI includes a **"Copy AI Prompt"** button. Click it and paste the prompt into your AI assistant when working on another project -- it will generate the code to make that service register itself with this registry on startup (fire-and-forget, failures silently ignored).

## License

[MIT](LICENSE)
