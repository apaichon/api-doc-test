package logger

import (
	"api/internal/db"
	"api/pkg/data"
	"api/pkg/data/models"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// ContactRepo represents the repository for contact operations
type SqliteLogger struct {
	DB *data.DB
}

type SqliteApiLogger struct {
	DB *db.DB
}

// NewContactRepo creates a new instance of ContactRepo
func NewSqliteLogger() *SqliteLogger {
	data := data.NewDB()
	return &SqliteLogger{DB: data}
}

func NewSqliteApiLogger() *SqliteApiLogger {
	data := db.NewDB()
	return &SqliteApiLogger{DB: data}
}

func (logger *SqliteApiLogger) InsertApiLog(logEntries []ApiLog) error {
	/*query := `
	  INSERT INTO api_logs (
	      level,
	      request_id,
	      timestamp,
	      method,
	      path,
	      status_code,
	      status_text,
	      duration,
	      request_body,
	      client_ip,
	      client_browser,
	      client_browser_version,
	      client_os,
	      client_os_version,
	      client_device,
	      user_id,
	      error

	  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	  `*/
	query := `
	INSERT INTO api_logs (log) VALUES (?)
	`

	// Prepare the SQL statement
	stmt, err := logger.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()

	// Iterate through the log entries and insert each one
	/*for _, logEntry := range logEntries {
		_, err := stmt.Exec(
			logEntry.Level,
			logEntry.RequestID,
			logEntry.Timestamp,
			logEntry.Method,
			logEntry.Path,
			logEntry.StatusCode,
			logEntry.StatusText,
			logEntry.Duration.Seconds(),
			logEntry.RequestBody,
			logEntry.ClientIP,
			logEntry.ClientBrowser,
			logEntry.ClientBrowserVersion,
			logEntry.ClientOS,
			logEntry.ClientOSVersion,
			logEntry.ClientDevice,
			logEntry.Error,
		)

		if err != nil {
			return fmt.Errorf("error inserting log: %w", err)
		}
	}*/

	for _, logEntry := range logEntries {
		// Convert logEntry to JSON string
		logEntryJSON, err := json.Marshal(logEntry)
		if err != nil {
			return fmt.Errorf("error marshaling log entry to JSON: %w", err)
		}

		_, err = stmt.Exec(string(logEntryJSON)) // Use the JSON string
		if err != nil {
			return fmt.Errorf("error inserting log: %w", err)
		}
	}

	return nil

}

// Insert multiple LogModel entries into the SQLite database
func (logger *SqliteLogger) InsertLog(logEntries []models.LogModel) error {
	// Prepare the SQL insert statement
	query := `
    INSERT INTO logs (
        log_id,
        timestamp,
        user_id,
        action,
        resource,
        status,
        client_ip,
        client_device,
        client_os,
        client_os_ver,
        client_browser,
        client_browser_ver,
        duration,
        errors
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	// Prepare the SQL statement
	stmt, err := logger.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()

	// Iterate through the log entries and insert each one
	for _, logEntry := range logEntries {
		_, err := stmt.Exec(
			logEntry.LogId,
			logEntry.Timestamp,
			logEntry.UserId,
			logEntry.Actions,
			logEntry.Resource,
			logEntry.Status,
			logEntry.ClientIp,
			logEntry.ClientDevice,
			logEntry.ClientOs,
			logEntry.ClientOsVersion,
			logEntry.ClientBrowser,
			logEntry.ClientBrowserVersion,
			logEntry.Duration.Nanoseconds(),
			logEntry.Errors,
		)

		if err != nil {
			return fmt.Errorf("error inserting log: %w", err)
		}
	}

	return nil
}
