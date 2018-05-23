package jsongofpdf

import (
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cast"
)

func (p *JSONGOFPDF) PreMultiCellFormField(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	attribute := p.GetString("attribute", logic, "")
	target := p.GetString("target", logic, "")
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0)

	field := p.Parser.FieldRegistry[p.RowIndex]
	if target != "" {
		for _, v := range p.Parser.FieldRegistry {
			if v.PathString == target || v.Key == target {
				field = v
			}
		}
	}

	if field.Key != "" {
		text := ""
		switch attribute {
		case "title":
			text = p.tr(strings.Replace(field.Title, "<br>", "\n", -1))
			break
		case "value":
			text = p.tr(strings.Replace(cast.ToString(field.Value), "<br>", "\n", -1))
		}

		cellList := pdf.SplitLines([]byte(text), width)
		cellHeight := float64(len(cellList)) * height
		if cellHeight > p.RowHeight {
			p.RowHeight = cellHeight
		}
	}

	return pdf
}

// PreRow calculates the size of the multi cell field before rendering it in the row so all rows are of equal height
func (p *JSONGOFPDF) PreRowMultiCell(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	attribute := p.GetString("attribute", logic, "")
	target := p.GetString("target", logic, "")
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0) // Line height of each cell, not cell height
	text := p.GetString("text", logic, "")

	cell := p.Tables[p.TableIndex].Rows[p.RowIndex].Cells[p.CellIndex]
	for _, row := range p.Tables[p.TableIndex].Rows {
		for _, rowCell := range row.Cells {
			if rowCell.Key == target || rowCell.Path == target {
				cell = rowCell
			}
		}
	}

	if text != "" {
		cell = Cell{
			Value: text,
		}
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
	if cellHeight > p.RowHeight {
		p.RowHeight = cellHeight
	}

	return pdf
}
