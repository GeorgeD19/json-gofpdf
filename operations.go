package jsongofpdf

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cast"
)

// UpdateX sets current pdfX position to currentX position plus what is passed.
func (p *JSONGOFPDF) UpdateX(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetX(p.CurrentX + p.GetFloat("width", logic, 0.0))
	p.CurrentX = pdf.GetX()
	return pdf
}

func (p *JSONGOFPDF) UpdateY(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	height := p.GetFloat("height", logic, 0.0)
	p.ManualY = pdf.GetY() + height
	pdf.SetY(p.ManualY)
	return pdf
}

// TableFunc uses json to render out a table using the passed data in the options.
func (p *JSONGOFPDF) TableFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	p.TableIndex = p.GetInt("index", logic, 0)
	pdf = p.Body(pdf, p.GetString("body", logic, ""))
	p.TableIndex = 0
	return pdf
}

// Body iterates over the table rows and renders the table based on the passed operations.
func (p *JSONGOFPDF) Body(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	// We store the logic of each row logic so we can alternate between each
	rowLogic := make([]string, 0)
	jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			// Operations that will be performed on every cell in this row
			operations := p.GetString("row", string(value), "")
			if operations != "" {
				rowLogic = append(rowLogic, operations)
			}
			break
		}
	})

	// Then foreach table.rows we can alternate between each function using RowIndex which resets on each new table.row
	if len(rowLogic) > 0 {
		rowLength := len(p.Tables[p.TableIndex].Rows)
		for x := 0; x < rowLength; x++ {

			p.RowIndex = x
			p.RowHeight = 0
			p.RowCells = 0.0
			p.CellIndex = 0
			p.CellPreIndex = 0

			for y := 0; y < len(p.Tables[p.TableIndex].Rows[x].Cells); y++ {
				p.PreOperations(pdf, rowLogic[p.RowFuncIndex])
				p.CellPreIndex++
			}

			pdf = p.RunArrayOperations(pdf, rowLogic[p.RowFuncIndex])

			if (p.RowFuncIndex + 1) < len(rowLogic) {
				p.RowFuncIndex++
			} else {
				p.RowFuncIndex = 0
			}
			if p.NextY > p.CurrentY {
				p.CurrentY = p.NextY
			}
			p.CurrentRowY = pdf.GetY()
		}
	}

	return pdf
}

func (p *JSONGOFPDF) Image(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	src := p.GetString("src", logic, "")
	name := p.GetString("name", logic, "")
	x := p.GetFloat("x", logic, 0.0)
	y := p.GetFloat("y", logic, 0.0)
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0)
	flow := p.GetBool("flow", logic, false)
	link := p.GetInt("link", logic, -1)
	linkStr := p.GetString("linkstr", logic, "")

	image, err := GetImage(src)
	if err == nil {
		imageContent := strings.TrimPrefix(image.Data, "0x")
		imageDecoded, err := hex.DecodeString(imageContent)
		if err == nil {
			options := gofpdf.ImageOptions{
				ReadDpi:   false,
				ImageType: image.Type,
			}
			pdf.RegisterImageOptionsReader(name, options, bytes.NewReader(imageDecoded))
			pdf.ImageOptions(name, x, y, width, height, flow, options, link, linkStr)
		}
	}

	return pdf
}

// SetXCurrentX sets pdfX to CurrentX position
func (p *JSONGOFPDF) SetXCurrentX(pdf *gofpdf.Fpdf) (opdf *gofpdf.Fpdf) {
	pdf.SetX(p.CurrentX)
	return pdf
}

// RowY sets pdf Y to CurrentRowY position
func (p *JSONGOFPDF) RowY(pdf *gofpdf.Fpdf) (opdf *gofpdf.Fpdf) {
	pdf.SetY(p.CurrentRowY)
	return pdf
}

// Line creates a line
func (p *JSONGOFPDF) Line(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	x := p.GetFloat("x", logic, 0.0)
	y := p.GetFloat("y", logic, 0.0)
	auto := p.GetString("auto", logic, "")
	if auto != "" {
		switch auto {
		case "P":
			y = p.CurrentY
			break
		case "C":
			y = pdf.GetY()
			break
		case "M":
			y = p.ManualY
			break
		}
	}
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 1.0)
	pdf.Line(x, y, x+width, y+height)
	return pdf
}

// LineRow creates a line at CurrentRowY position. Pass "width", "height" float object properties.
func (p *JSONGOFPDF) LineRow(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	currentX := pdf.GetX()
	rowY := p.CurrentRowY
	nextX := currentX + p.GetFloat("width", logic, 0.0)
	nextY := rowY + p.GetFloat("height", logic, 0.0)
	pdf.Line(currentX, rowY, nextX, nextY)
	return pdf
}

func (p *JSONGOFPDF) SetInitY(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetY(p.initY)
	return pdf
}

func (p *JSONGOFPDF) MultiCell(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	attribute := p.GetString("attribute", logic, "")
	target := p.GetString("target", logic, "")
	loop := p.GetBool("loop", logic, false)

	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0) // Line height of each cell, not cell height
	border := p.GetString("border", logic, "")
	align := p.GetString("align", logic, "L")
	text := p.GetString("text", logic, "")
	fill := p.GetBool("fill", logic, false)
	format := p.GetString("format", logic, "")

	if v := p.GetString("calculation", logic, ""); v != "" {
		text = p.Calculation(v, text)
	}

	cell := Cell{}
	if target == "" {
		cell = p.Tables[p.TableIndex].Rows[p.RowIndex].Cells[p.CellIndex]
	}

	for _, rowCell := range p.Tables[p.TableIndex].Rows[p.RowIndex].Cells {
		if rowCell.Key == target || rowCell.Path == target {
			cell = rowCell
		}
	}

	if loop {
		for _, table := range p.Tables {
			for _, row := range table.Rows {
				for _, rowCell := range row.Cells {
					if rowCell.Key == target || rowCell.Path == target {
						cell = rowCell
					}
				}
			}
		}
	}

	// fmt.Println(cell)

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

	cellCount := 0.0
	cellX := pdf.GetX()

	renderText := ""
	// if cast.ToString(cell.Value) != "" {
	switch attribute {
	case "title":
		renderText = cell.Title
		break
	case "value":
		renderText = cast.ToString(cell.Value)
		break
	default:
		renderText = text
	}

	if cell.Type == "currency" || cell.Type == "date" {
		format = cell.Type
	}

	if format != "" {
		renderText = p.Format(format, renderText)
	}

	renderText = p.tr(strings.Replace(renderText, "<br>", "\n", -1))
	cellList := pdf.SplitLines([]byte(renderText), width)
	cellCount = float64(len(cellList))

	if renderText != "" {
		pdf.MultiCell(width, height, renderText, border, align, fill)
	}

	// }

	if cellCount < p.RowCells {
		for i := 0; i < int(p.RowCells-cellCount); i++ {
			pdf.SetX(cellX)
			pdf.MultiCell(width, height, "", border, align, fill)
		}
	}

	if cell.Images != nil && attribute == "value" {
		for _, image := range cell.Images {

			// For any media against the field
			imageDecoded, err := hex.DecodeString(image.Data)
			if err == nil {
				options := gofpdf.ImageOptions{
					ReadDpi:   false,
					ImageType: image.Type,
				}

				imageWidth := float64(image.Width) / float64(p.DPI)
				imageHeight := float64(image.Height) / float64(p.DPI)

				if imageWidth > width {
					imageWidth = width
					imageHeight = 0
				}

				// Pass binary image into PDF
				pdf.RegisterImageOptionsReader(string(p.MediaIndex), options, bytes.NewReader(imageDecoded))
				pdf.ImageOptions(string(p.MediaIndex), cellX, pdf.GetY(), imageWidth, imageHeight, true, options, -1, "")
				p.MediaIndex++
			} else {
				fmt.Println(err)
			}
		}
	}

	CurrentY := pdf.GetY()
	if CurrentY > p.NextY {
		p.NextY = CurrentY
	}

	return pdf
}
