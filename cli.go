package unicreds

import (
	"encoding/csv"
	"io"

	"github.com/olekukonko/tablewriter"
)

const (
	// TableFormatTerm format the table for a terminal session
	TableFormatTerm = iota // 0
	// TableFormatCSV format the table as CSV
	TableFormatCSV // 1
)

// TableWriter enables writing of tables in a variety of formats
type TableWriter struct {
	tableFormat int
	headers     []string
	rows        [][]string
	wr          io.Writer
}

// NewTable create a new table writer
func NewTable(wr io.Writer) *TableWriter {
	return &TableWriter{wr: wr}
}

// SetHeaders set the column headers
func (tw *TableWriter) SetHeaders(headers []string) {
	tw.headers = headers
}

// SetFormat set the format
func (tw *TableWriter) SetFormat(tableFormat int) {
	tw.tableFormat = tableFormat
}

func (tw *TableWriter) Write(row []string) {
	tw.rows = append(tw.rows, row)
}

// BulkWrite append an array of rows to the buffer
func (tw *TableWriter) BulkWrite(rows [][]string) {
	tw.rows = append(tw.rows, rows...)
}

// Render render the table out to the supplied writer
func (tw *TableWriter) Render() error {
	switch tw.tableFormat {
	case TableFormatTerm:
		table := tablewriter.NewWriter(tw.wr)

		// Convert headers to interface slice
		headerInterfaces := make([]interface{}, len(tw.headers))
		for i, h := range tw.headers {
			headerInterfaces[i] = h
		}
		table.Header(headerInterfaces...)

		// Convert rows to interface slices
		for _, row := range tw.rows {
			rowInterfaces := make([]interface{}, len(row))
			for i, cell := range row {
				rowInterfaces[i] = cell
			}
			table.Append(rowInterfaces...)
		}
		table.Render()
	case TableFormatCSV:
		w := csv.NewWriter(tw.wr)

		for _, r := range tw.rows {
			if err := w.Write(r); err != nil {
				return err
			}
		}
		w.Flush()

		if err := w.Error(); err != nil {
			return err
		}
	}

	return nil
}
