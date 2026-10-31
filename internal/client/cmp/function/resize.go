package function

import "fyne.io/fyne/v2/widget"

func SetDefaultColumnsWidthCard(table *widget.Table) {
	colWidths := []float32{150, 150, 150, 200, 150, 50, 150, 150, 150, 150}
	for idx, colWidth := range colWidths {
		table.SetColumnWidth(idx, colWidth)
	}
}

func SetDefaultColumnsWidthText(table *widget.Table) {
	colWidths := []float32{150, 150, 500, 150, 150, 150}
	for idx, colWidth := range colWidths {
		table.SetColumnWidth(idx, colWidth)
	}
}

func SetDefaultColumnsWidthLoginPassword(table *widget.Table) {
	colWidths := []float32{150, 150, 150, 150, 150, 150, 150}
	for idx, colWidth := range colWidths {
		table.SetColumnWidth(idx, colWidth)
	}
}

func SetDefaultColumnsWidthBinary(table *widget.Table) {
	colWidths := []float32{150, 150}
	for idx, colWidth := range colWidths {
		table.SetColumnWidth(idx, colWidth)
	}
}
