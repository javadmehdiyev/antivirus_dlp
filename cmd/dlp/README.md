# DLP Agent

Agent for testing Data Loss Prevention (DLP) systems. Checks if DLP blocks transmission of sensitive data.

## Build

```bash
go build -o dlp ./cmd/dlp
```

## Usage

```bash
# Run with go run
go run cmd/dlp/main.go [-file <path>]... [-url <url>] [-method <HTTP_METHOD>] [-json <json_file>]

# Run compiled binary
./dlp [-file <path>]... [-url <url>] [-method <HTTP_METHOD>] [-json <json_file>]

# Run without specifying files (uses default test files)
./dlp
```

## Parameters

- `-file` - Path to test file (optional, can be used multiple times to test multiple files)
- `-url` - Target URL for DLP check (optional, defaults to URL from settings API)
- `-method` - HTTP method: GET, POST, PUT, etc. (default: GET)
- `-json` - Path to JSON file to store results (default: `dlp_results.json`)

## Default Behavior

If no files are specified with `-file`, the agent will:
1. Check if the following 4 default test files exist:
   - `test_credit_card.txt` (category: `credit_card`)
   - `test_passport.txt` (category: `passport_number`)
   - `test_dlp_data.csv` (category: `file_upload_csv`)
   - `test_dlp_data.xlsx` (category: `file_upload_xlsx`)
2. If all files exist, use them for testing
3. If any file is missing, create all 4 files automatically
4. Process all files sequentially, sending one request per file

## Configuration

If `-url` is not provided, the agent automatically retrieves the DLP URL from the settings API:
- Settings endpoint: `http://127.0.0.1:8000/api/settings-agent`
- The agent fetches the `url_dlp` from the settings response

## Examples

```bash
# Test with default files (auto-creates if missing)
./dlp

# Test with default files and custom URL
./dlp -url https://testdlp.net/

# Test single file
./dlp -file test_dlp_data.txt

# Test multiple files
./dlp -file test1.txt -file test2.txt -file test3.txt

# Test with POST method
./dlp -file test_dlp_data.txt -url https://testdlp.net/ -method POST

# Custom JSON output file
./dlp -file test_dlp_data.txt -url https://testdlp.net/ -json custom_results.json
```

## Test Files

The agent supports testing with files containing sensitive data:
- Credit card numbers
- Social Security Numbers (SSN)
- Passport numbers
- Personal information
- API keys, passwords
- CSV files with sensitive data
- Excel (XLSX) files with sensitive data

### Default Test Files

When no files are specified, the agent uses 4 default test files with different categories:

1. **test_credit_card.txt** - Contains credit card numbers (category: `credit_card`)
2. **test_passport.txt** - Contains passport information (category: `passport_number`)
3. **test_dlp_data.csv** - CSV file with sensitive data (category: `file_upload_csv`)
4. **test_dlp_data.xlsx** - Excel file with sensitive data (category: `file_upload_xlsx`)

Each file is processed sequentially, and results are saved with file name and category information.

## Output

When processing files, you'll see progress indicators:
```
[1/4] Processing file: test_credit_card.txt
DLP Active: false
Status: Request succeeded: 200 OK

[2/4] Processing file: test_passport.txt
...
```

- `[N/M] Processing file: <filename>` - Progress indicator showing current file number
- `DLP Active: false` - DLP did not block the request, file sent successfully
- `DLP Active: true` - DLP blocked the request
- `Status: <message>` - Detailed status message
- Exit code `0` - All files processed successfully, no DLP detected
- Exit code `1` - DLP detected in at least one file

## How It Works

1. If no files are specified, checks for default test files or creates them
2. For each file:
   - Reads the file content
   - Sends HTTP request with file content as multipart form-data to the specified URL
   - Checks if the request was blocked (error = blocked)
   - Saves result with file name and category to JSON
3. Processes files sequentially (one request per file)
4. Returns result indicating if DLP is active in any file

## Results Storage

Results are saved to a JSON file (default: `dlp_results.json`) with the following structure:

```json
{
  "results": [
    {
      "timestamp": "2025-12-10T16:56:39.262418+04:00",
      "status_text": "Request succeeded: 200 OK",
      "is_dlp_active": false,
      "file_name": "test_credit_card.txt",
      "category": "credit_card"
    }
  ]
}
```

Each entry includes:
- `timestamp` - Timestamp of the check
- `status_text` - Detailed status message
- `is_dlp_active` - Whether DLP blocked the request
- `file_name` - Name of the processed file
- `category` - Category of the file (`credit_card`, `passport_number`, `file_upload_csv`, `file_upload_xlsx`)

The JSON file keeps only the last 15 entries.





