package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVExporter struct{}

func (CSVExporter) Export(path string, rows []TotalRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{"student_id", "name", "total_percent"})
	for _, r := range rows {
		w.Write([]string{r.StudentID, r.Name, fmt.Sprintf("%.2f", r.Total)})
	}
	w.Flush()
	return w.Error()
}
