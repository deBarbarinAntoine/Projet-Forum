package utils

import (
	"Projet-Forum/internal/models"
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type Logs []Log

type Log struct {
	Time       time.Time      `json:"time"`
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	ReqId      int            `json:"req_id,omitempty"`
	User       models.Session `json:"user,omitempty"`
	ClientIP   string         `json:"client_ip,omitempty"`
	ReqMethod  string         `json:"req_method,omitempty"`
	ReqURL     string         `json:"req_url,omitempty"`
	HttpStatus int            `json:"http_status,omitempty"`
	ErrOutput  string         `json:"output,omitempty"`
}

var Logger *slog.Logger
var logs *os.File
var wg sync.WaitGroup

// closeLog
//
//	@Description: closes the log file.
func closeLog() {
	if logs != nil {
		err := logs.Close()
		if err != nil {
			log.Println(GetCurrentFuncName(), slog.Any("output", err))
		}
	}
}

// LogInit is meant to be run as a goroutine to create a new log file every day
// appending the file's creation timestamp in its name.
func LogInit() {
	duration := SetDailyTimer(0)
	var jsonHandler *slog.JSONHandler
	var err error
	var filename string
	defer closeLog()

	// checking if logs directory exists
	_, err = os.Stat(Path + "logs")
	if os.IsNotExist(err) {

		// create it if it doesn't exist
		err = os.Mkdir(Path+"logs", 0750)
		if err != nil {
			log.Println("LogInit(): error when creating directory 'logs'", err)
		}
	}

	for {
		filename = Path + "logs/logs_" + time.Now().Format(time.DateOnly) + ".log"
		closeLog()
		logs, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(GetCurrentFuncName(), slog.Any("output", err))
		}
		jsonHandler = slog.NewJSONHandler(logs, nil)
		Logger = slog.New(jsonHandler)
		Logger.Info(GetCurrentFuncName(), slog.String("goroutine", "LogInit"))
		time.Sleep(duration)
		duration = time.Hour * 24
	}
}

// fetchLogInfo retrieves all Log from `file` and stores it in *log.
func (log *Logs) fetchLogInfo(file string) {
	defer wg.Done()
	filename := "logs/" + file
	data, err := os.ReadFile(filename)
	if len(data) == 0 {
		return
	}
	lines := bytes.Split(data, []byte("\n"))
	var singleLog Log
	for _, line := range lines {
		err = json.Unmarshal(line, &singleLog)
		if err != nil {
			return
		}
		*log = append(*log, singleLog)
	}
}

// printFileNames
//
//	@Description: verbose function, for testing
//	@param files
//	@return []string
//	@return for
func printFileNames(files []os.DirEntry) []string {
	var result []string
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result
}

// RetrieveLogs fetches all Log from all files *.log in /logs directory
// and returns a Logs array.
func RetrieveLogs() (logArray Logs) {
	logFiles, err := os.ReadDir(Path + "logs/.")
	//fmt.Printf("logFiles: %#v\n", printFileNames(logFiles)) // verbose
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	} else {
		reg := regexp.MustCompile(`\.log$`)
		for _, file := range logFiles {
			if reg.MatchString(file.Name()) {
				wg.Add(1)
				go logArray.fetchLogInfo(file.Name())
			}
		}
	}
	wg.Wait()
	logArray.sortLogs()
	return logArray
}

// sortLogs sort all Log from the newest to the oldest.
func (log *Logs) sortLogs() {
	sort.Slice(*log, func(i, j int) bool {
		return (*log)[i].Time.After((*log)[j].Time)
	})
}

// FetchAttrLogs filters Log returning only Log matching the given `level`.
func FetchAttrLogs(attr string, value string) Logs {
	attr = strings.ToLower(attr)
	logs := RetrieveLogs()
	var result Logs
	switch attr {
	case "level":
		switch strings.ToUpper(value) {
		case "INFO", "WARN", "ERROR":
			for _, singleLog := range logs {
				if singleLog.Level == strings.ToUpper(value) {
					result = append(result, singleLog)
				}
			}
		default:
			return nil
		}
	case "user", "username":
		for _, singleLog := range logs {
			if strings.EqualFold(singleLog.User.Username, value) {
				result = append(result, singleLog)
			}
		}
	default:
		return nil
	}
	result.sortLogs()
	return result
}
