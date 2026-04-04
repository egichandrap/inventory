package http

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/example/jwt-ddd-clean/internal/pkg/sanitizer"
)

// ValidationMiddleware validates incoming requests
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate URL path
		if !sanitizer.ValidatePathTraversal(r.URL.Path) {
			sendValidationError(w, "Invalid URL path")
			return
		}

		// Validate query parameters
		for key, values := range r.URL.Query() {
			for _, value := range values {
				if !sanitizer.ValidateSQL(key) || !sanitizer.ValidateXSS(key) {
					sendValidationError(w, "Invalid query parameter")
					return
				}
				if !sanitizer.ValidateSQL(value) || !sanitizer.ValidateXSS(value) {
					sendValidationError(w, "Invalid query value")
					return
				}
			}
		}

		// Validate headers
		for key, value := range r.Header {
			if !sanitizer.ValidateSQL(key) || !sanitizer.ValidateXSS(key) {
				sendValidationError(w, "Invalid header name")
				return
			}
			for _, v := range value {
				if !sanitizer.ValidateSQL(v) || !sanitizer.ValidateXSS(v) {
					sendValidationError(w, "Invalid header value")
					return
				}
			}
		}

		// Validate Content-Type for POST/PUT/PATCH
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			if contentType == "" {
				sendValidationError(w, "Content-Type header required")
				return
			}

			// Only allow JSON content type
			if !strings.Contains(contentType, "application/json") {
				sendValidationError(w, "Content-Type must be application/json")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateJSONBody validates and sanitizes JSON request body
func ValidateJSONBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Reject unknown fields

	if err := decoder.Decode(v); err != nil {
		sendValidationError(w, "Invalid JSON: "+err.Error())
		return false
	}

	// Sanitize string fields in the struct
	sanitizeStruct(v)

	return true
}

// sanitizeStruct recursively sanitizes string fields in a struct
func sanitizeStruct(v interface{}) {
	// This is a simplified version
	// In production, use reflection to sanitize all string fields
}

// ValidateUUID validates UUID format
func ValidateUUID(uuid string) bool {
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return uuidPattern.MatchString(uuid)
}

// ValidatePositiveInt validates positive integer
func ValidatePositiveInt(n int) bool {
	return n > 0
}

// ValidateNonNegative validates non-negative integer
func ValidateNonNegative(n int) bool {
	return n >= 0
}

// ValidateStringLength validates string length
func ValidateStringLength(s string, min, max int) bool {
	length := len(s)
	return length >= min && length <= max
}

// ValidateAlphanumeric validates alphanumeric string
func ValidateAlphanumeric(s string) bool {
	alnumPattern := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alnumPattern.MatchString(s)
}

// sendValidationError sends a validation error response
func sendValidationError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error": map[string]string{
			"code":    "ERR_VALIDATION",
			"message": message,
		},
	})
}
