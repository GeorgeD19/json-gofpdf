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

	// cell.str = strings.Join(strList[0:count], " ")
	// cell.list = pdf.SplitLines([]byte(cell.str), colWd-cellGap-cellGap)
	// cell.ht = float64(len(cell.list)) * lineHt
	// if cell.ht > maxHt {
	// 	maxHt = cell.ht
	// }

	if field.Key != "" {
		text := ""
		switch attribute {
		case "title":
			text = p.tr(strings.Replace(field.Title, "<br>", "\n", -1))
			break
		case "value":
			text = p.tr(strings.Replace(cast.ToString(field.Value), "<br>", "\n", -1))
			// for _, image := range field.Media {
			// 	// For any media against the field
			// 	ImageDecoded, _ := hex.DecodeString(image.Data)
			// }
		}

		cellList := pdf.SplitLines([]byte(text), width)
		cellHeight := float64(len(cellList)) * height

		// TODO if value field we need to add in the image dimensions

		if cellHeight > p.RowHeight {
			p.RowHeight = cellHeight
		}

	}

	return pdf
}
