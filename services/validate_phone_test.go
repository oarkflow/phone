package services

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePhonePreservesOrderAndTruncatesOutput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.csv")
	output := filepath.Join(dir, "output.csv")
	if err := os.WriteFile(input, []byte("name,phone\nAlice,+14155552671\nBob,invalid\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(output, []byte("stale bytes that must disappear"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := ValidatePhone(input, output, "phone", ',', ','); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(output)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 || records[1][0] != "Alice" || records[2][0] != "Bob" {
		t.Fatalf("unexpected records: %v", records)
	}
	if records[1][2] != "+14155552671" || records[1][3] != "false" {
		t.Fatalf("unexpected valid enrichment: %v", records[1])
	}
	if records[2][3] != "true" {
		t.Fatalf("unexpected invalid enrichment: %v", records[2])
	}
}

func TestValidatePhoneRequiresPhoneColumn(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.csv")
	if err := os.WriteFile(input, []byte("name\nAlice\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePhone(input, filepath.Join(dir, "output.csv"), "phone", ',', ','); err == nil {
		t.Fatal("expected missing phone column error")
	}
}
