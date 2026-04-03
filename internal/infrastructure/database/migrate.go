package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations runs all pending migrations using direct SQL execution
// This approach works for both PostgreSQL and SQLite
func RunMigrations(db *DB) error {
	// Get migrations directory path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	execDir := filepath.Dir(execPath)
	migrationsDir := filepath.Join(execDir, "migrations")

	// Fallback to working directory
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		wd, _ := os.Getwd()
		migrationsDir = filepath.Join(wd, "migrations")
	}

	// Read migration files
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("failed to glob migration files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationsDir)
	}

	// Sort files by name to ensure correct order
	sort.Strings(files)

	fmt.Printf("📁 Found %d migration file(s)\n", len(files))

	// Execute each migration
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filepath.Base(file), err)
		}

		// Extract UP migration content
		sql := extractUpMigration(string(content))
		if sql == "" {
			fmt.Printf("⏭️  Skipping %s (no UP migration)\n", filepath.Base(file))
			continue
		}

		fmt.Printf("🔄 Running migration: %s\n", filepath.Base(file))

		// Execute SQL statements (split by semicolon)
		statements := splitSQLStatements(sql)
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			// Skip pure comment-only statements
			if strings.HasPrefix(stmt, "--") && !strings.Contains(stmt, ";") {
				continue
			}

			// Remove single-line comments from the beginning of the statement
			for strings.HasPrefix(stmt, "--") {
				lines := strings.SplitN(stmt, "\n", 2)
				if len(lines) > 1 {
					stmt = strings.TrimSpace(lines[1])
				} else {
					break
				}
			}

			if stmt == "" {
				continue
			}

			if _, err := db.Exec(stmt); err != nil {
				// Ignore certain PostgreSQL errors (e.g., already exists)
				if isPostgresIgnoreError(err) {
					fmt.Printf("⚠️  Warning in %s: %v\n", filepath.Base(file), err)
					continue
				}
				return fmt.Errorf("failed to execute migration %s: %w\nSQL: %s", filepath.Base(file), err, stmt)
			}
		}
	}

	fmt.Println("✅ Database migrations completed successfully")
	return nil
}

// isPostgresIgnoreError returns true if the error should be ignored for PostgreSQL
func isPostgresIgnoreError(err error) bool {
	msg := err.Error()
	// Ignore "already exists" errors
	if strings.Contains(msg, "already exists") {
		return true
	}
	if strings.Contains(msg, "duplicate key") {
		return true
	}
	return false
}

// extractUpMigration extracts the UP migration content from a SQL file
func extractUpMigration(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inUpSection := false

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "-- +migrate Up") {
			inUpSection = true
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "-- +migrate Down") {
			inUpSection = false
			continue
		}
		if inUpSection {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// splitSQLStatements splits SQL content by semicolons
func splitSQLStatements(sql string) []string {
	var statements []string
	current := strings.Builder{}

	for _, char := range sql {
		current.WriteRune(char)
		if char == ';' {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}
