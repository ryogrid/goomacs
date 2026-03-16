package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestComprehensiveZZZ_Summary runs last (alphabetically after all other
// TestComprehensive* tests) and prints a summary of the test suite results
// from the failure log.
func TestComprehensiveZZZ_Summary(t *testing.T) {
	data, err := os.ReadFile(failureLogPath)
	if err != nil {
		t.Logf("No failure log found: %v", err)
		printSummary(t, nil)
		return
	}

	var records []FailureRecord
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("Failed to parse failure log: %v", err)
	}

	printSummary(t, records)

	// Fail the summary test if there are critical or high severity failures
	for _, r := range records {
		if r.Severity == "critical" || r.Severity == "high" {
			t.Errorf("Found %s severity failure: %s - %s", r.Severity, r.TestCaseID, r.TestName)
		}
	}
}

func printSummary(t *testing.T, records []FailureRecord) {
	t.Helper()

	bySeverity := map[string]int{}
	byCategory := map[string]int{}
	for _, r := range records {
		bySeverity[r.Severity]++
		byCategory[r.Category]++
	}

	total := len(records)
	fmt.Println("\n========================================")
	fmt.Println("  Comprehensive E2E Test Suite Summary")
	fmt.Println("========================================")
	fmt.Printf("  Total failures recorded: %d\n", total)
	if total > 0 {
		fmt.Println("\n  Failures by severity:")
		for _, sev := range []string{"critical", "high", "medium", "low"} {
			if count, ok := bySeverity[sev]; ok {
				fmt.Printf("    %-10s: %d\n", sev, count)
			}
		}
		fmt.Println("\n  Failures by category:")
		for cat, count := range byCategory {
			fmt.Printf("    %-20s: %d\n", cat, count)
		}
	} else {
		fmt.Println("  All tests passed!")
	}
	fmt.Println("========================================")
}
