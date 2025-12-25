# Combined Antivirus & DLP Agent

A combined agent that runs Antivirus and DLP checks sequentially in a single build.

## Building

```bash
go build -o agent ./cmd/combined
```

Or for cross-platform builds:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o agent-linux ./cmd/combined

# Windows
GOOS=windows GOARCH=amd64 go build -o agent.exe ./cmd/combined

# macOS
GOOS=darwin GOARCH=amd64 go build -o agent-macos ./cmd/combined
```

## Usage

### Basic usage (run both checks)

```bash
./agent
```

### With DLP files specified

```bash
./agent -file test1.txt -file test2.txt
```

### With DLP URL specified

```bash
./agent -dlp-url http://example.com/api/dlp
```

### Skip Antivirus check

```bash
./agent -skip-antivirus
```

### Skip DLP check

```bash
./agent -skip-dlp
```

### All options

```bash
./agent \
  -file test1.txt \
  -file test2.txt \
  -antivirus-json antivirus_results.json \
  -dlp-json dlp_results.json \
  -dlp-url http://example.com/api/dlp \
  -method POST \
  -skip-antivirus \
  -skip-dlp
```

## Flags

- `-file`: Path to file for DLP check (can be used multiple times)
- `-antivirus-json`: Path to JSON file for saving Antivirus results (default: `antivirus_results.json`)
- `-dlp-json`: Path to JSON file for saving DLP results (default: `dlp_results.json`)
- `-dlp-url`: URL for DLP check (if not specified, will be obtained from settings)
- `-method`: HTTP method for DLP requests (default: `GET`)
- `-skip-antivirus`: Skip Antivirus check
- `-skip-dlp`: Skip DLP check

## Exit Codes

- `0`: All checks passed successfully
- `1`: Problem detected (virus or DLP active)

## Execution Order

1. **Step 1**: Antivirus check is started
2. **Step 2**: DLP check is started (if Antivirus passed successfully)
3. **Result**: Overall result of all checks is displayed

If a problem is detected at any stage, the program will exit with error code 1.


