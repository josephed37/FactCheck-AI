// backend/internal/database/database.go

package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/josephed37/FactCheck-AI/backend/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// DB is a global variable to hold the database connection pool.
var DB *sql.DB

// InitDB initializes the SQLite database connection and creates the necessary tables.
func InitDB(dataSourceName string) error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS fact_checks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		statement TEXT NOT NULL,
		verdict TEXT NOT NULL,
		confidence TEXT NOT NULL,
		reason TEXT,
		additional_context TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create fact_checks table: %w", err)
	}

	DB = db
	return nil
}

// SaveFactCheck inserts a new fact-check record into the database.
func SaveFactCheck(req models.FactCheckRequest, resp models.GeminiResponse) error {
	insertSQL := `INSERT INTO fact_checks(statement, verdict, confidence, reason, additional_context) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(insertSQL, req.Statement, resp.Verdict, resp.Confidence, resp.Reason, resp.AdditionalContext)
	if err != nil {
		return fmt.Errorf("failed to insert fact-check record: %w", err)
	}
	return nil
}

// NEW: GetFactCheckHistory retrieves all records from the fact_checks table.
// It returns a slice of FactCheckHistoryItem structs, ordered by the most recent first.
func GetFactCheckHistory() ([]models.FactCheckHistoryItem, error) {
	// 1. Define the SQL query to select all columns from our table.
	//    We order by 'created_at DESC' to get the newest entries first.
	querySQL := `SELECT id, statement, verdict, confidence, reason, additional_context, created_at FROM fact_checks ORDER BY created_at DESC`

	// 2. Execute the query. `DB.Query` is used for SELECT statements that
	//    can return multiple rows.
	rows, err := DB.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("failed to execute history query: %w", err)
	}
	// `defer rows.Close()` is crucial. It ensures that the connection to the
	// database is closed when the function finishes, preventing resource leaks.
	defer rows.Close()

	// 3. Prepare a slice to hold our results.
	// A slice is Go's version of a dynamic array or list.
	var history []models.FactCheckHistoryItem

	// 4. Iterate over the results.
	// `rows.Next()` advances to the next row in the result set. The loop
	// continues as long as there are more rows to process.
	for rows.Next() {
		// For each row, create a temporary variable to scan the data into.
		var item models.FactCheckHistoryItem

		// 5. Scan the data from the current row into our 'item' struct.
		// The order of the arguments to `rows.Scan` must exactly match the
		// order of the columns in our SELECT statement.
		if err := rows.Scan(&item.ID, &item.Statement, &item.Verdict, &item.Confidence, &item.Reason, &item.AdditionalContext, &item.CreatedAt); err != nil {
			// If scanning fails for any row, we return an error.
			return nil, fmt.Errorf("failed to scan history row: %w", err)
		}

		// 6. Append the successfully scanned item to our history slice.
		history = append(history, item)
	}

	// 7. After the loop, return the complete history slice.
	return history, nil
}
