package csvexport

import (
	"encoding/csv"
	"fmt"
	"io"
	"pos/internal/application"
	"strconv"
)

func Write(output io.Writer, results []application.StudentResult) error {
	// UTF-8 BOM helps spreadsheet applications detect the encoding.
	if _, err := io.WriteString(output, "\uFEFF"); err != nil {
		return fmt.Errorf("write UTF-8 BOM: %w", err)
	}
	writer := csv.NewWriter(output)
	if err := writer.Write([]string{"student_id", "name", "group", "final_score"}); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}
	for _, result := range results {
		record := []string{string(result.StudentID), result.Name, result.Group, strconv.Itoa(result.Score)}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write CSV result: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush CSV: %w", err)
	}
	return nil
}
