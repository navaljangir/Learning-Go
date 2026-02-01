# Hot Reload in Go (Air)

## What is Hot Reload?

Hot reload automatically restarts your server when you save a file - no manual restart needed.

## Node.js vs Go

| Node.js | Go |
|---------|-----|
| `nodemon` | `air` |
| `nodemon.json` | `.air.toml` |
| `npx nodemon server.js` | `air` |

## Installing Air

```bash
# Install globally (one time)
go install github.com/air-verse/air@latest
```

### Important: Add Go bin to PATH

After installing, you need to add Go's bin folder to your system PATH (one time setup).

**Where is it installed?**
```bash
# Check your GOPATH
go env GOPATH

# Air is at: <GOPATH>/bin/air
# Example: C:\Users\jangi\go\bin\air.exe
```

**Windows (PowerShell - temporary for current session):**
```powershell
$env:Path += ";$(go env GOPATH)\bin"
air  # Now works!
```

**Windows (Permanent - recommended):**
1. Press `Win + R`, type `sysdm.cpl`, press Enter
2. Go to **Advanced** tab → **Environment Variables**
3. Under **User variables**, find `Path` and click **Edit**
4. Click **New** and add: `%USERPROFILE%\go\bin`
5. Click OK → OK → OK
6. **Restart your terminal** (close and reopen VS Code)

**macOS/Linux (Permanent):**
```bash
# Add to ~/.bashrc or ~/.zshrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

**Verify it works:**
```bash
air --version
```

## Using Air

```bash
# Navigate to your project
cd my_project/

# Option 1: Generate config and run
air init    # Creates .air.toml (optional - air works without it too)
air         # Start watching

# Option 2: Run with defaults (no config needed)
air
```

## What Happens

```
$ air
  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \_

watching...
building...
running...

# Now edit any .go file and save
# Air automatically rebuilds and restarts!
```

## .air.toml Configuration

```toml
# Auto-generated with: air init

root = "."              # Project root
tmp_dir = "tmp"         # Temp directory for builds

[build]
cmd = "go build -o ./tmp/main.exe ."   # Build command
bin = "./tmp/main.exe"                  # Binary to run
include_ext = ["go", "html", "tmpl"]    # Watch these extensions
exclude_dir = ["tmp", "vendor"]         # Ignore these folders
delay = 1000                            # Wait 1s before rebuild (ms)

[log]
time = true             # Show timestamps

[misc]
clean_on_exit = true    # Delete tmp/ when air stops
```

## Common Configuration Changes

### Run a specific file or folder
```toml
[build]
# Run specific file (if you only have one main file)
cmd = "go build -o ./tmp/main.exe app.go"

# Run from a subfolder
cmd = "go build -o ./tmp/main.exe ./cmd/api"

# Default: compile entire current package
cmd = "go build -o ./tmp/main.exe ."
```

### Watch more file types
```toml
include_ext = ["go", "html", "css", "js", "json", "yaml"]
```

### Exclude test files from triggering rebuild
```toml
exclude_regex = ["_test.go"]
```

### Change build output location
```toml
cmd = "go build -o ./bin/server ."
bin = "./bin/server"
```

## Comparison: Running Without vs With Air

### Without Air (manual restart)
```bash
# Terminal 1
go run .
# Make changes...
# Ctrl+C to stop
go run .    # Manually restart every time
```

### With Air (automatic)
```bash
air
# Make changes...
# Server restarts automatically!
```

## Troubleshooting

### "air: command not found"
```bash
# Add Go bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Or on Windows, add to System Environment Variables:
# %USERPROFILE%\go\bin
```

### Air not detecting changes
```bash
# Make sure you're in the right directory
# Check .air.toml exclude_dir isn't blocking your files
```

## Quick Reference

| Command | Description |
|---------|-------------|
| `go install github.com/air-verse/air@latest` | Install Air |
| `air init` | Generate .air.toml config |
| `air` | Start watching & running |
| `air -c .air.custom.toml` | Use custom config file |
