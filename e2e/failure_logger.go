package e2e

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"
)

// FailureRecord represents a single test failure entry.
type FailureRecord struct {
	TestCaseID    string `json:"test_case_id"`
	TestName      string `json:"test_name"`
	InputSequence string `json:"input_sequence"`
	ExpectedState string `json:"expected_state"`
	ActualState   string `json:"actual_state"`
	Severity      string `json:"severity"`
	Category      string `json:"category"`
	Timestamp     string `json:"timestamp"`
}

var (
	failureMu      sync.Mutex
	failureLogPath = "tasks/e2e-failures.json"
)

// ResetFailureLog clears the failure log file, starting fresh for a new test run.
func ResetFailureLog() error {
	failureMu.Lock()
	defer failureMu.Unlock()
	return os.WriteFile(failureLogPath, []byte("[]"), 0644)
}

// RecordFailure appends a failure record to the failure log file.
// Severity levels: "critical" (crash/data loss), "high" (wrong behavior),
// "medium" (cosmetic), "low" (minor).
func RecordFailure(t *testing.T, caseID string, input string, expected string, actual string, severity string, category string) {
	t.Helper()

	record := FailureRecord{
		TestCaseID:    caseID,
		TestName:      t.Name(),
		InputSequence: input,
		ExpectedState: expected,
		ActualState:   actual,
		Severity:      severity,
		Category:      category,
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	failureMu.Lock()
	defer failureMu.Unlock()

	// Read existing records
	var records []FailureRecord
	data, err := os.ReadFile(failureLogPath)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &records)
	}

	// Append new record
	records = append(records, record)

	// Write back
	out, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		t.Errorf("failed to marshal failure records: %v", err)
		return
	}

	if err := os.WriteFile(failureLogPath, out, 0644); err != nil {
		t.Errorf("failed to write failure log: %v", err)
	}
}
