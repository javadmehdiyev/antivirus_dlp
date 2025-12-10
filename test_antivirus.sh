#!/bin/bash
# Test shell script for antivirus scanning
# Safe test file with patterns that might trigger checks

# Suspicious patterns (but safe)
ENCODED=$(echo "test string" | base64)
DECODED=$(echo "$ENCODED" | base64 -d)

# Network operations (often checked)
curl -s http://example.com/test > /dev/null

# File operations (monitored)
if [ -f "/tmp/test.txt" ]; then
    cat /tmp/test.txt
fi

# Process operations
ps aux | grep test

echo "This is a safe test script"





