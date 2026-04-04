package sanitizer

import (
	"html"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// sqlPattern matches common SQL injection patterns
	sqlPattern = regexp.MustCompile(`(?i)(\b(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|ALTER|CREATE|EXEC|EXECUTE)\b|--|;|/\*|\*/|@@|@|CHAR|NCHAR|VARCHAR|NVARCHAR|TRIM|CONVERT|CAST|DECLARE|SET|TABLE|DATABASE|FROM|WHERE|HAVING|GROUP BY|ORDER BY)`)
	
	// xssPattern matches common XSS patterns
	xssPattern = regexp.MustCompile(`(?i)(<script|javascript:|on\w+\s*=|<iframe|<object|<embed|<link|<style|<img\s+.*\son\w+=)`)
	
	// pathTraversalPattern matches path traversal attempts
	pathTraversalPattern = regexp.MustCompile(`(\.\./|\.\.\\|%2e%2e%2f|%2e%2e/|\.%2e/|%2e\.\/)`)
)

// SanitizeString sanitizes a string input
func SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// Trim whitespace
	input = strings.TrimSpace(input)

	// HTML escape
	input = html.EscapeString(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	return input
}

// SanitizeHTML sanitizes HTML input
func SanitizeHTML(input string) string {
	return html.EscapeString(input)
}

// ValidateSQL checks for SQL injection attempts
func ValidateSQL(input string) bool {
	if input == "" {
		return true
	}
	return !sqlPattern.MatchString(input)
}

// ValidateXSS checks for XSS attempts
func ValidateXSS(input string) bool {
	if input == "" {
		return true
	}
	return !xssPattern.MatchString(input)
}

// ValidatePathTraversal checks for path traversal attempts
func ValidatePathTraversal(input string) bool {
	if input == "" {
		return true
	}
	return !pathTraversalPattern.MatchString(input)
}

// SanitizeFilename sanitizes file names
func SanitizeFilename(filename string) string {
	// Remove path components
	filename = filepath.Base(filename)
	
	// Remove dangerous characters
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.ReplaceAll(filename, "\x00", "")
	
	return filename
}

// ValidateEmail validates and sanitizes email
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	
	// Basic email validation
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

// ValidateUsername validates username format
func ValidateUsername(username string) bool {
	if username == "" {
		return false
	}
	
	// Username: 3-50 chars, alphanumeric + underscore only
	usernamePattern := regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)
	return usernamePattern.MatchString(username)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password minimal 8 karakter"
	}
	
	if len(password) > 128 {
		return false, "Password maksimal 128 karakter"
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	
	if !hasUpper {
		return false, "Password harus mengandung huruf kapital"
	}
	
	if !hasLower {
		return false, "Password harus mengandung huruf kecil"
	}
	
	if !hasNumber {
		return false, "Password harus mengandung angka"
	}
	
	if !hasSpecial {
		return false, "Password harus mengandung karakter khusus"
	}
	
	return true, ""
}

// IsCommonPassword checks if password is in common passwords list
func IsCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "123456", "12345678", "qwerty", "abc123",
		"monkey", "1234567", "letmein", "trustno1", "dragon",
		"baseball", "iloveyou", "master", "sunshine", "ashley",
		"bailey", "shadow", "123456789", "654321", "superman",
		"qazwsx", "michael", "football", "password1", "password123",
	}
	
	passwordLower := strings.ToLower(password)
	for _, common := range commonPasswords {
		if passwordLower == common {
			return true
		}
	}
	
	return false
}

// SanitizeSearchQuery sanitizes search query
func SanitizeSearchQuery(query string) string {
	// Remove SQL injection patterns
	query = sqlPattern.ReplaceAllString(query, "")
	
	// Remove XSS patterns
	query = xssPattern.ReplaceAllString(query, "")
	
	// Trim and normalize whitespace
	query = strings.TrimSpace(query)
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")
	
	return query
}
