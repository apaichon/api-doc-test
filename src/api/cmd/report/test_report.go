package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type TestEvent struct {
	Action  string  `json:"Action"`
	Test    string  `json:"Test"`
	Output  string  `json:"Output"`
	Package string  `json:"Package"`
	Time    string  `json:"Time"`
	Elapsed float64 `json:"Elapsed"`
}

type TestReport struct {
	TotalTests  int     `json:"total"`
	PassedTests int     `json:"passed"`
	FailedTests int     `json:"failed"`
	Coverage    float64 `json:"coverage"`
	Duration    string  `json:"duration"`
}

func generateReport() error {
	cmd := exec.Command("go", "test",
		"-json", "../server/...",
		"-coverprofile=coverage.out")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run tests: %v", err)
	}

	report := parseTestOutput(output)
	jsonReport, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}

	fmt.Println(string(jsonReport))
	return nil
}

func parseTestOutput(output []byte) *TestReport {
	var report TestReport
	var startTime time.Time
	var endTime time.Time

	// Split output into lines and parse each test event
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		var event TestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		switch event.Action {
		case "start":
			if startTime.IsZero() {
				startTime, _ = time.Parse(time.RFC3339Nano, event.Time)
			}
			report.TotalTests++
		case "pass":
			report.PassedTests++
			endTime, _ = time.Parse(time.RFC3339Nano, event.Time)
		case "fail":
			report.FailedTests++
			endTime, _ = time.Parse(time.RFC3339Nano, event.Time)
		}
	}

	// Calculate duration
	if !startTime.IsZero() && !endTime.IsZero() {
		report.Duration = endTime.Sub(startTime).String()
	}

	// TODO: Parse coverage from coverage.out file
	report.Coverage = 0.0

	return &report
}
