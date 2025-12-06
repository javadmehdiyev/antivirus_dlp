# Test PowerShell script for antivirus scanning
# Safe test file with patterns that antivirus might check

# Suspicious patterns (but safe)
$encoded = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes("test"))
$decoded = [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($encoded))

# Network operations (often monitored)
$webClient = New-Object System.Net.WebClient
$url = "http://example.com/test"

# File operations (monitored by antivirus)
$testPath = "C:\temp\test.txt"
if (Test-Path $testPath) {
    Get-Content $testPath
}

Write-Host "This is a safe test script"


