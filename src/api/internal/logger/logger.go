package logger

import (
	"api/config"
	"api/pkg/data/models"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type ApiLog struct {
	Level                string
	RequestID            string
	Timestamp            time.Time
	Method               string
	Path                 string
	StatusCode           int
	StatusText           string
	Duration             time.Duration
	RequestBody          string
	ClientIP             string
	ClientBrowser        string
	ClientBrowserVersion string
	ClientOS             string
	ClientOSVersion      string
	ClientDevice         string
	UserID               int
	Error                string
}

type Logger struct {
	once        sync.Once
	currentFile *os.File
	currentPath string
	logMutex    sync.Mutex
}

// FileWithModTime is a struct to hold file information along with its modification time
type FileWithModTime struct {
	Name    string
	ModTime time.Time
	IsDir   bool
}

// The singleton instance of LogInitializer
var loggerInstance *Logger
var loggerOnce sync.Once
var relativePath string = "../../logs"

// GetLogInitializer returns the singleton instance of LogInitializer
func GetLogInitializer() *Logger {
	loggerOnce.Do(func() {
		loggerInstance = &Logger{}
	})
	return loggerInstance
}

// Initialize ensures that setup logic is done only once
func (li *Logger) Initialize() {
	li.once.Do(func() {
		relativePath := "../../logs"
		if err := os.MkdirAll(relativePath, 0755); err != nil {
			log.Fatalf("Error creating logs directory: %v", err)
		}
	})
}

// WriteLogToFile writes the log entry to a file, handling initialization and file management
func (li *Logger) WriteLogToFile(logEntry models.LogModel) {
	// Ensure initialization happens only once
	li.Initialize()
	conf := config.GetConfig()

	// Determine the log file name based on the current timestamp (every 5 minutes)
	logFileName := fmt.Sprintf("%04d-%02d-%02d-%02d_%02d.log",
		logEntry.Timestamp.Year(),
		logEntry.Timestamp.Month(),
		logEntry.Timestamp.Day(),
		logEntry.Timestamp.Hour(),
		(logEntry.Timestamp.Minute()/conf.LogMergeMin)*conf.LogMergeMin,
	)

	// relativePath := "../../logs"
	logFilePath := filepath.Join(relativePath, logFileName)
	absolutePath, err := filepath.Abs(logFilePath)
	if err != nil {
		log.Fatalf("Error obtaining absolute path: %v", err)
	}

	// Use a mutex to ensure thread safety
	li.logMutex.Lock()
	defer li.logMutex.Unlock()

	// If the current file path changes, close the previous file and open a new one
	if li.currentPath != absolutePath {
		if li.currentFile != nil {
			li.currentFile.Close()
		}

		li.currentPath = absolutePath

		// Open the new file
		li.currentFile, err = os.OpenFile(absolutePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
	}

	// Write the log entry as JSON to the file
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Fatalf("Error marshaling log data: %v", err)
	}

	_, err = li.currentFile.Write(append(jsonData, '\n')) // Add newline
	if err != nil {
		log.Fatalf("Error writing to log file: %v", err)
	}
}

func (li *Logger) WriteApiLogToFile(logEntry ApiLog) {
	// Ensure initialization happens only once
	li.Initialize()
	conf := config.GetConfig()

	// Determine the log file name based on the current timestamp (every 5 minutes)
	// time := logEntry.Timestamp.Format("2025-01-12-18_39")
	logFileName := fmt.Sprintf("%04d-%02d-%02d-%02d_%02d.log",
		logEntry.Timestamp.Year(),
		logEntry.Timestamp.Month(),
		logEntry.Timestamp.Day(),
		logEntry.Timestamp.Hour(),
		(logEntry.Timestamp.Minute()/conf.LogMergeMin)*conf.LogMergeMin,
	)
	// logFileName := GetPeriodFileName(logEntry.Timestamp, conf.LogMergeMin)

	fmt.Printf("logFileName: %v, conf.LogMergeMin: %v\n", logFileName, conf.LogMergeMin)

	// relativePath := "../../logs"
	logFilePath := filepath.Join(relativePath, logFileName)
	absolutePath, err := filepath.Abs(logFilePath)
	if err != nil {
		log.Printf("Error obtaining absolute path: %v", err)
		return
	}

	// Use a mutex to ensure thread safety
	li.logMutex.Lock()
	defer li.logMutex.Unlock()

	// If the current file path changes, close the previous file and open a new one
	if li.currentPath != absolutePath {
		if li.currentFile != nil {
			li.currentFile.Close()
		}

		li.currentPath = absolutePath

		// Open the new file
		li.currentFile, err = os.OpenFile(absolutePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening log file: %v", err)
		}
	}

	// Write the log entry as JSON to the file
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log data: %v", err)
		return
	}

	_, err = li.currentFile.Write(append(jsonData, '\n')) // Add newline
	if err != nil {
		log.Printf("Error writing to log file: %v", err)
		return
	}
}

func GetPeriodFileName(timestamp time.Time, periodMinutes int) string {
	// Calculate the start of the period
	minutes := timestamp.Minute()
	periodStart := minutes - (minutes % periodMinutes)

	// Create new time with period start
	fileTime := timestamp.Truncate(time.Hour).Add(time.Duration(periodStart) * time.Minute)

	// Format filename with period start time
	return fmt.Sprintf("%s.log", fileTime.Format("2006-01-02_15_04"))
}

func (li *Logger) MoveApiLogsToSQLite() {

	absolutePath, err := filepath.Abs(relativePath)
	
	if err != nil {
		log.Printf("Error reading logs directory: %v", err)
		return
	}

	// Get and sort the files by the oldest modification time
	files, err := listFilesOrderedByOldest(absolutePath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	if len(files) == 0 {
		return // No logs to process
	}

	// Read the log file and insert into SQLite
	logFilePath := filepath.Join(absolutePath, files[0].Name)
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var logList []ApiLog
	for scanner.Scan() {
		var logEntry ApiLog
		if err := json.Unmarshal([]byte(scanner.Text()), &logEntry); err != nil {
			log.Printf("Error unmarshaling log data: %v", err)
			continue
		}

		log.Printf("inserting into SQLite: %v", logEntry)
		logList = append(logList, logEntry)
		// fmt.Printf("%v", logList)
	}

	logger := NewSqliteApiLogger()
	logger.InsertApiLog(logList)

	// Delete the log file after processing
	err = os.Remove(logFilePath)
	if err != nil {
		log.Printf("Error deleting log file: %v", err)
	}

}

// Function to read the last log file and insert its content into SQLite
func (li *Logger) MoveLogsToSQLite() {

	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		log.Printf("Error reading logs directory: %v", err)
		return
	}

	// Get and sort the files by the oldest modification time
	files, err := listFilesOrderedByOldest(absolutePath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	if len(files) == 0 {
		return // No logs to process
	}

	// Read the log file and insert into SQLite
	logFilePath := filepath.Join(absolutePath, files[0].Name)
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var logList []models.LogModel
	for scanner.Scan() {
		var logEntry models.LogModel
		if err := json.Unmarshal([]byte(scanner.Text()), &logEntry); err != nil {
			log.Printf("Error unmarshaling log data: %v", err)
			continue
		}

		log.Printf("inserting into SQLite: %v", logEntry)
		logList = append(logList, logEntry)
		// fmt.Printf("%v", logList)
	}

	logger := NewSqliteLogger()
	logger.InsertLog(logList)

	// Delete the log file after processing
	err = os.Remove(logFilePath)
	if err != nil {
		log.Printf("Error deleting log file: %v", err)
	}

}

// Function to read the last log file and insert its content into SQLite
func (li *Logger) MoveLogsToPostgres() {

	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		log.Printf("Error reading logs directory: %v", err)
		return
	}

	// Get and sort the files by the oldest modification time
	files, err := listFilesOrderedByOldest(absolutePath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	if len(files) == 0 {
		return // No logs to process
	}

	// Read the log file and insert into SQLite
	logFilePath := filepath.Join(absolutePath, files[0].Name)
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var logList []models.LogModel
	for scanner.Scan() {
		var logEntry models.LogModel
		if err := json.Unmarshal([]byte(scanner.Text()), &logEntry); err != nil {
			log.Printf("Error unmarshaling log data: %v", err)
			continue
		}

		log.Printf("inserting into Postgres: %v", logEntry)
		logList = append(logList, logEntry)
		// fmt.Printf("%v", logList)
	}

	logger := NewPostgresLogger()
	logger.InsertLog(logList)

	// Delete the log file after processing
	err = os.Remove(logFilePath)
	if err != nil {
		log.Printf("Error deleting log file: %v", err)
	}

}

// Read and sort files by modification time in ascending order (oldest to newest)
func listFilesOrderedByOldest(path string) ([]FileWithModTime, error) {
	// Read the directory
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileInfos []FileWithModTime

	// Get file information and modification time
	for _, file := range files {
		info, err := file.Info() // Get detailed information about the file
		if err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, FileWithModTime{
			Name:    file.Name(),
			ModTime: info.ModTime(), // Use ModTime as a proxy for creation time
			IsDir:   file.IsDir(),
		})
	}

	// Sort files by modification time
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime.Before(fileInfos[j].ModTime) // Ascending order
	})

	return fileInfos, nil
}
