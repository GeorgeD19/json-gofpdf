package jsongofpdf

import (
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cast"
)

// PreOperations will iterate through the array of operations and execute each
func (p *JSONGOFPDF) PreOperations(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			pdf = p.PreParseObject(pdf, string(value))
			break
		}
	})

	return pdf
}

// PreParseObject entry point
func (p *JSONGOFPDF) PreParseObject(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		pdf = p.PreRunOperation(pdf, string(key), string(value))
		return nil
	})

	return pdf
}

// PreRunOperation determines which function to run for the pre-render pipeline
func (p *JSONGOFPDF) PreRunOperation(pdf *gofpdf.Fpdf, name string, logic string) (opdf *gofpdf.Fpdf) {
	switch name {
	case "multicell":
		p.PreRowMultiCell(pdf, logic)
		break
	}
	return pdf
}

// PreRowMultiCell calculates the size of the multi cell field before rendering it in the row so all rows are of equal height
func (p *JSONGOFPDF) PreRowMultiCell(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	attribute := p.GetString("attribute", logic, "")
	target := p.GetString("target", logic, "")
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0) // Line height of each cell, not cell height
	text := p.GetString("text", logic, "")

	cell := p.Tables[p.TableIndex].Rows[p.RowIndex].Cells[p.CellPreIndex]

	if text != "" {
		cell = Cell{
			Value: text,
		}
	}

	if target != "" {
		if cell.Path != target {
			cell.Disabled = true
		}
		for index, value := range p.Globals {
			if index == target {
				cell = Cell{
					Path:  index,
					Key:   index,
					Title: index,
					Value: value,
				}
			}
		}
	}

	if cell.Disabled == false {
		renderText := ""
		switch attribute {
		case "title":
			renderText = p.tr(strings.Replace(cell.Title, "<br>", "\n", -1))
			break
		case "value":
			renderText = p.tr(strings.Replace(cast.ToString(cell.Value), "<br>", "\n", -1))
		}

		cellList := pdf.SplitLines([]byte(renderText), width)
		cellHeight := float64(len(cellList)) * height
		cellCount := float64(len(cellList))

		if cellCount > p.RowCells {
			p.RowCells = cellCount
		}

		if cellHeight > p.RowHeight {
			p.RowHeight = cellHeight
		}
	}

	// This stops the first cell being left behind
	_, pageHeight := pdf.GetPageSize()
	if pdf.GetY()+p.RowHeight > pageHeight {
		pdf.AddPage()
	}

	return pdf
}
