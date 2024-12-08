package helpers

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
)

func CreateColumns(columns []string, rows [][]string, viewportWidth int) []table.Column {
	maxWidths := make([]int, len(columns))
	for i, col := range columns {
		maxWidths[i] = len(col)
	}

	for _, row := range rows {
		for i, cell := range row {
			cellLength := len(cell)
			if cellLength > maxWidths[i] {
				maxWidths[i] = cellLength
			}
		}
	}

	totalWidth := 0
	for _, width := range maxWidths {
		totalWidth += width
	}
	scaleFactor := float64(viewportWidth) / float64(totalWidth)

	cols := make([]table.Column, 0, len(columns))
	for i, title := range columns {
		cols = append(cols, table.Column{
			Title: strings.ToUpper(title),
			Width: int(float64(maxWidths[i]) * scaleFactor),
		})
	}

	return cols
}

func CreateRows(inputRows [][]string) []table.Row {
	rows := make([]table.Row, 0, len(inputRows))
	for _, row := range inputRows {
		rows = append(rows, row)
	}
	return rows
}
