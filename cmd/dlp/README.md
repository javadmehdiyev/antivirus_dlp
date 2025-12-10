# DLP Agent

Agent for testing Data Loss Prevention (DLP) systems. Checks if DLP blocks transmission of sensitive data.

## Build

```bash
go build -o dlp ./cmd/dlp
```

## Usage

```bash
# Run with go run
go run cmd/dlp/main.go -file <path> -url <url> [-method <HTTP_METHOD>] [-json <json_file>]

# Run compiled binary
./dlp -file <path> -url <url> [-method <HTTP_METHOD>] [-json <json_file>]
```

## Parameters

- `-file` - Path to test file (required, default: `test_dlp_data.txt`)
- `-url` - Target URL for DLP check (required, defaults to URL from settings API)
- `-method` - HTTP method: GET, POST, PUT, etc. (default: POST)
- `-json` - Path to JSON file to store results (default: `dlp_results.json`)

## Configuration

If `-url` is not provided, the agent automatically retrieves the DLP URL from the settings API:
- Settings endpoint: `http://127.0.0.1:8000/api/settings-agent`
- The agent fetches the `url_dlp` from the settings response

## Examples

```bash
# Test with default file and settings URL
./dlp -file test_dlp_data.txt

# Test with custom URL
./dlp -file test_dlp_data.txt -url https://testdlp.net/ -method POST

# GET request
./dlp -file test_dlp_data.txt -url https://testdlp.net/ -method GET

# Custom JSON output file
./dlp -file test_dlp_data.txt -url https://testdlp.net/ -json custom_results.json
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
- `Status: <message>` - Detailed status message
- Exit code `0` - Request succeeded
- Exit code `1` - Request blocked

## How It Works

1. Reads the test file content
2. Sends HTTP request with file content to the specified URL
3. Checks if the request was blocked (error = blocked)
4. Returns result indicating if DLP is active

## Results Storage

Results are saved to a JSON file (default: `dlp_results.json`) with the following structure:
- Timestamp of the check
- Status text
- DLP active status

The JSON file keeps only the last 15 entries.





