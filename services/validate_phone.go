package services

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/oarkflow/phone"
)

var enrichmentColumns = []string{
	"region",
	"phone_type_label",
	"carrier_name",
	"carrier_mnc",
	"carrier_mcc",
	"carrier_nnc",
	"phone_type_code",
	"dial_code",
}

// GetCsvHeader reads and normalizes the first CSV record. It is retained for
// compatibility; ValidatePhone uses encoding/csv directly so multiline fields
// and read errors are handled correctly.
func GetCsvHeader(scanner *bufio.Scanner, comma rune) map[int]string {
	columns := make(map[int]string)
	if !scanner.Scan() {
		return columns
	}
	reader := csv.NewReader(strings.NewReader(scanner.Text()))
	reader.Comma = comma
	header, err := reader.Read()
	if err != nil {
		return columns
	}
	for index, column := range header {
		columns[index] = clean([]byte(column))
	}
	return columns
}

type csvJob struct {
	index  int
	record []string
}

type csvResult struct {
	index  int
	record []string
}

// ValidatePhone streams csvFile, enriches each row in parallel, and writes
// results in input order. Existing output files are truncated at open, and all
// CSV read/write errors are returned to the caller.
func ValidatePhone(csvFile, out, phoneKey string, comma, outputComma rune) error {
	input, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("open input CSV: %w", err)
	}
	defer input.Close()

	output, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open output CSV: %w", err)
	}
	defer output.Close()

	reader := csv.NewReader(bufio.NewReader(input))
	reader.Comma = comma
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read CSV header: %w", err)
	}
	for index := range header {
		header[index] = clean([]byte(header[index]))
	}

	phoneIndex := -1
	for index, column := range header {
		if column == phoneKey {
			phoneIndex = index
			break
		}
	}
	if phoneIndex < 0 {
		return fmt.Errorf("phone column %q not found", phoneKey)
	}

	validatedColumn := "validated_" + phoneKey
	invalidColumn := "is_invalid_" + phoneKey
	outputHeader := append(append([]string(nil), header...), validatedColumn, invalidColumn)
	outputHeader = append(outputHeader, enrichmentColumns...)

	writer := csv.NewWriter(output)
	writer.Comma = outputComma
	if err := writer.Write(outputHeader); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}

	workerCount := max(1, runtime.GOMAXPROCS(0))
	jobs := make(chan csvJob, workerCount)
	results := make(chan csvResult, workerCount)
	readErrors := make(chan error, 1)
	var workers sync.WaitGroup
	workers.Add(workerCount)
	for range workerCount {
		go func() {
			defer workers.Done()
			for job := range jobs {
				number := phone.Verify(job.record[phoneIndex])
				record := append(append([]string(nil), job.record...), number.Phone, fmt.Sprint(number.Invalid))
				record = append(record,
					number.CountryCode,
					number.PhoneTypeHuman,
					number.CarrierName,
					number.CarrierMnc,
					number.CarrierMcc,
					number.CarrierNnc,
					fmt.Sprint(number.PhoneType),
					fmt.Sprint(number.DialCode),
				)
				results <- csvResult{index: job.index, record: record}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for index := 0; ; index++ {
			record, err := reader.Read()
			if errors.Is(err, io.EOF) {
				readErrors <- nil
				return
			}
			if err != nil {
				readErrors <- fmt.Errorf("read CSV record %d: %w", index+2, err)
				return
			}
			jobs <- csvJob{index: index, record: record}
		}
	}()
	go func() {
		workers.Wait()
		close(results)
	}()

	next := 0
	pending := make(map[int][]string, workerCount)
	for result := range results {
		pending[result.index] = result.record
		for {
			record, found := pending[next]
			if !found {
				break
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("write CSV record %d: %w", next+2, err)
			}
			delete(pending, next)
			next++
		}
	}
	if err := <-readErrors; err != nil {
		return err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush output CSV: %w", err)
	}
	return nil
}

func clean(value []byte) string {
	var normalized strings.Builder
	normalized.Grow(len(value))
	for _, char := range string(value) {
		if char == ' ' || char == '_' || char == '-' ||
			(char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') {
			normalized.WriteRune(char)
		}
	}
	return normalized.String()
}

// ProcessNumber enriches legacy line-oriented jobs. New code should prefer
// ValidatePhone, which handles quoted multiline CSV records and ordered output.
func ProcessNumber(jobs <-chan []byte, results chan<- map[string]string, wg *sync.WaitGroup, columns map[int]string, phoneKey string, comma rune) {
	defer wg.Done()
	for job := range jobs {
		reader := csv.NewReader(strings.NewReader(string(job)))
		reader.Comma = comma
		fields, err := reader.Read()
		if err != nil {
			continue
		}
		data := make(map[string]string, len(fields)+len(enrichmentColumns)+2)
		for index, value := range fields {
			if column, found := columns[index]; found {
				data[column] = value
			}
		}
		number := phone.Verify(data[phoneKey])
		data["validated_"+phoneKey] = number.Phone
		data["is_invalid_"+phoneKey] = fmt.Sprint(number.Invalid)
		data["region"] = number.CountryCode
		data["phone_type_label"] = number.PhoneTypeHuman
		data["carrier_name"] = number.CarrierName
		data["carrier_mnc"] = number.CarrierMnc
		data["carrier_mcc"] = number.CarrierMcc
		data["carrier_nnc"] = number.CarrierNnc
		data["phone_type_code"] = fmt.Sprint(number.PhoneType)
		data["dial_code"] = fmt.Sprint(number.DialCode)
		results <- data
	}
}
