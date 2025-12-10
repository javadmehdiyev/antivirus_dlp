// Test JavaScript file for antivirus scanning
// This file contains patterns that might trigger antivirus checks
// but is completely safe - no actual malware

function testFunction() {
    var x = "test";
    var y = eval("x"); // eval() is often flagged by antivirus
    console.log(y);
    
    // Suspicious patterns (but safe)
    var encoded = btoa("test string");
    var decoded = atob(encoded);
    
    return decoded;
}

// Network-related code (often checked by antivirus)
var url = "http://example.com/test";
fetch(url).then(response => response.text());





