# DLP Agent

Agent for testing Data Loss Prevention (DLP) systems. Checks if DLP blocks transmission of sensitive data.

## Build

```bash
go build -o dlp ./cmd/dlp
```

## Usage

```bash
# Run with go run
go run cmd/dlp/main.go -file <path> -url <url> [-method <HTTP_METHOD>]

# Run compiled binary
./dlp -file <path> -url <url> [-method <HTTP_METHOD>]
```

## Parameters

- `-file` - Path to test file (required)
- `-url` - Target URL for DLP check (required)
- `-method` - HTTP method: GET, POST, PUT, etc. (default: POST)

## Examples

```bash
# Test with testdlp.net
./dlp -file test_dlp_data.txt -url https://testdlp.net/ -method POST

# GET request
./dlp -file test_dlp_data.txt -url https://testdlp.net/ -method GET
```

## Test Files

Use files containing sensitive data:
- Credit card numbers
- Social Security Numbers (SSN)
- Personal information
- API keys, passwords

Example: `test_dlp_data.txt`

## Output

- `DLP Active: false` - DLP did not block the request, file sent successfully
- `DLP Active: true` - DLP blocked the request
- Exit code `0` - Request succeeded
- Exit code `1` - Request blocked

## How It Works

1. Reads the test file content
2. Sends HTTP request with file content to the specified URL
3. Checks if the request was blocked (error = blocked)
4. Returns result indicating if DLP is active

