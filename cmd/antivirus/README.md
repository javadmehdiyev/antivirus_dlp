# Antivirus Agent

Agent for testing antivirus systems. Checks if antivirus detects viruses or malware in downloaded files.

## Build

```bash
go build -o antivirus ./cmd/antivirus
```

## Usage

```bash
# Run with go run
go run cmd/antivirus/main.go [-json <json_file>]

# Run compiled binary
./antivirus [-json <json_file>]
```

## Parameters

- `-json` - Path to JSON file to store results (default: `antivirus_results.json`)

## Configuration

The agent automatically retrieves the antivirus service URL from the settings API:
- Settings endpoint: `http://127.0.0.1:8000/api/settings-agent`
- The agent fetches the `url_antivirus` from the settings response

## Examples

```bash
# Run with default JSON output file
./antivirus

# Run with custom JSON output file
./antivirus -json custom_results.json
```

## How It Works

1. Retrieves antivirus service URL from settings API (`http://127.0.0.1:8000/api/settings-agent`)
2. Sends GET request to download a file from the antivirus service endpoint
3. If download succeeds, saves the file to `uploads/` directory
4. Waits 5 seconds and checks if the file still exists
5. If file was deleted, antivirus detected a virus
6. Returns result indicating if virus was detected

## Output

- `Virus Detected: false` - Antivirus did not detect virus, file downloaded and still exists
- `Virus Detected: true` - Antivirus detected virus and deleted the file
- `File Name: <name>` - Name of the downloaded file (if available)
- `File Path: <path>` - Path where file was saved (if available)
- `File Exists: <true/false>` - Whether file still exists after 5 seconds
- `Status: <message>` - Detailed status message
- Exit code `0` - No virus detected
- Exit code `1` - Virus detected

## Results Storage

Results are saved to a JSON file (default: `antivirus_results.json`) with the following structure:
- Timestamp of the check
- File name and path
- Status text
- Virus detection result
- File existence status

The JSON file keeps only the last 15 entries.





