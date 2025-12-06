# Antivirus Agent

Agent for testing antivirus systems. Checks if antivirus detects viruses or malware in files.

## Build

```bash
go build -o antivirus ./cmd/antivirus
```

## Usage

```bash
# Run with go run
go run cmd/antivirus/main.go -file <path> -url <url> [-method <HTTP_METHOD>]

# Run compiled binary
./antivirus -file <path> -url <url> [-method <HTTP_METHOD>]
```

## Parameters

- `-file` - Path to test file (required)
- `-url` - Antivirus service URL (required)
- `-method` - HTTP method: GET, POST, PUT, etc. (default: POST)

## Examples

```bash
# Test with JavaScript file
./antivirus -file test_antivirus.js -url https://antivirus-service.com/scan -method POST

# Test with PowerShell script
./antivirus -file test_antivirus.ps1 -url https://antivirus-service.com/scan -method POST

# Test with shell script
./antivirus -file test_antivirus.sh -url https://antivirus-service.com/scan -method POST
```

## Test Files

Use files that might contain suspicious patterns:
- Executable files (.exe, .dll)
- Scripts (.js, .ps1, .sh)
- Archives (.zip, .rar)
- Files with encoded content

Examples:
- `test_antivirus.js` - JavaScript with suspicious patterns
- `test_antivirus.ps1` - PowerShell script
- `test_antivirus.sh` - Bash script
- `test_antivirus_data.txt` - EICAR-like pattern

## Output

- `Virus Detected: false` - Antivirus did not detect virus, file sent successfully
- `Virus Detected: true` - Antivirus detected virus and blocked the request
- Exit code `0` - Request succeeded
- Exit code `1` - Virus detected

## How It Works

1. Reads the test file content
2. Sends HTTP request with file content to antivirus service
3. Checks if the request was blocked (error = virus detected)
4. Returns result indicating if virus was detected




