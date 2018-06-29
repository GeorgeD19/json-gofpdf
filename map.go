package jsongofpdf

import (
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cast"
)

// Page sizes
// pt = 1.0
// mm 72.0 / 25.4
// cm 72.0 / 2.54
// in, "inch" 72.0

// a3 (w:841.89, h:1190.55) - NC
// a4 (w:210.00155555555557, h:297) // docWidth - 30 (margin from left / right) = 180.00155555555557
// a4L (w:297, h:210.00155555555557) // docWidth - 30 = 267
// a5 (w:420.94, h:595.28) - NC
// a6 (w:297.64, h:420.94) - NC
// a2 (w:1190.55, h:1683.78) - NC
// a1 (w:1683.78, h:2383.94) - NC
// letter (w:612, h:792) - NC
// legal (w:612, h:1008) - NC
// tabloid (w:792, h:1224) - NC

// mm = 72.0 / 25.4
// a3 (w:841.89, h:1190.55)
// a4 (w:8.267777777777778, h:33.145275590551181)

// SetHeaderFunc maps json to gofpdf SetHeaderFunc function.
func (p *JSONGOFPDF) SetHeaderFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		pdf = p.RunArrayOperations(pdf, logic)
		p.CurrentRowY = pdf.GetY()
		p.CurrentY = pdf.GetY()
	})

	return pdf
}

// New passes the orientation, unit, size and dir object properties to the gofpdf New function creating a new pdf
func (p *JSONGOFPDF) New(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	orientation := p.GetString("orientation", logic, "P")
	unit := p.GetString("unit", logic, "mm")
	size := p.GetString("size", logic, "A4")
	directory := p.GetString("dir", logic, "")
	return gofpdf.New(orientation, unit, size, directory)
}

// SetCellMargin maps json to gofpdf SetCellMargin function.
// Default is "margin": 0.0
func (p *JSONGOFPDF) SetCellMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetCellMargin(p.GetFloat("margin", logic, 0.0))
	return pdf
}

// SetLeftMargin maps json to gofpdf SetLeftMargin function. Pass "margin" as a float object property in json logic.
// Default is "margin": 0.0
func (p *JSONGOFPDF) SetLeftMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetLeftMargin(p.GetFloat("margin", logic, 0.0))
	return pdf
}

// SetMargins maps json to gofpdf SetMargins function. Pass "left", "top" and "right" as float object properties in json logic.
// Defaults are "left": 0.0, "top": 0.0, "right": 0.0
func (p *JSONGOFPDF) SetMargins(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetMargins(p.GetFloat("left", logic, 0.0), p.GetFloat("top", logic, 0.0), p.GetFloat("right", logic, 0.0))
	return pdf
}

// SetRightMargin maps json to gofpdf SetRightMargin function. Pass "margin" as float object property in json logic.
// Default is "margin": 0.0
func (p *JSONGOFPDF) SetRightMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetRightMargin(p.GetFloat("margin", logic, 0.0))
	return pdf
}

// SetTopMargin maps json to gofpdf SetTopMargin function. Pass "margin" as float object property in json logic.
// Default is "margin": 0.0
func (p *JSONGOFPDF) SetTopMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetTopMargin(p.GetFloat("margin", logic, 0.0))
	return pdf
}

// SetAutoPageBreak maps json to gofpdf SetAutoPageBreak function. Pass "auto" boolean and "margin" float object properties in json logic.
// Default is "auto": true, "margin": 15.0
func (p *JSONGOFPDF) SetAutoPageBreak(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetAutoPageBreak(p.GetBool("auto", logic, true), p.GetFloat("margin", logic, 15.0))
	return pdf
}

// SetDisplayMode maps json to gofpdf SetDisplayMode function. Pass "zoom" and "layout" string object properties in json logic.
// Default is "zoom": "", "layout": ""
func (p *JSONGOFPDF) SetDisplayMode(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetDisplayMode(p.GetString("zoom", logic, ""), p.GetString("layout", logic, ""))
	return pdf
}

// Ln maps json to gofpdf Ln function. Pass "height" as string object property in json logic.
// Default is "height": -1.0
func (p *JSONGOFPDF) Ln(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.Ln(p.GetFloat("height", logic, -1.0))
	return pdf
}

// SetDrawColor maps json to gofpdf SetDrawColor function. Pass "r", "g" and "b" as integer object properties in json logic.
// Defaults are "r": 0, "g": 0, "b": 0
func (p *JSONGOFPDF) SetDrawColor(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(p.GetInt("r", logic, 0), p.GetInt("g", logic, 0), p.GetInt("b", logic, 0))
	return pdf
}

// SetFillColor maps json to gofpdf SetFillColor function. Pass "r", "g" and "b" as integer object properties in json logic.
// Defaults are "r": 0, "g": 0, "b": 0
func (p *JSONGOFPDF) SetFillColor(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetFillColor(p.GetInt("r", logic, 0), p.GetInt("g", logic, 0), p.GetInt("b", logic, 0))
	return pdf
}

// SetTextColor maps json to gofpdf SetTextColor function. Pass "r", "g" and "b" as integer object properties in json logic.
// Defaults are "r": 0, "g": 0, "b": 0
func (p *JSONGOFPDF) SetTextColor(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetTextColor(p.GetInt("r", logic, 0), p.GetInt("g", logic, 0), p.GetInt("b", logic, 0))
	return pdf
}

// AddPage maps json to gofpdf AddPage function. No arguemenets are taken.
func (p *JSONGOFPDF) AddPage(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.AddPage()
	return pdf
}

// SetFont maps json to gofpdf SetFont function. Pass in "family" string, "style" string, "size" float properties in json logic.
// Defaults are "family": "Arial", "style": "", "size", 8.0
func (p *JSONGOFPDF) SetFont(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetFont(p.GetString("family", logic, "Arial"), p.GetString("style", logic, ""), p.GetFloat("size", logic, 8.0))
	return pdf
}

// AliasNbPages maps json to gofpdf AliasNbPages function. Pass in "alias" as string object property in json logic.
// Default is "alias": ""
func (p *JSONGOFPDF) AliasNbPages(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.AliasNbPages(p.GetString("alias", logic, ""))
	return pdf
}

// SetY maps json to gofpdf SetY function. Pass "y" as float object property in json logic.
// Default is "y": 0.0
func (p *JSONGOFPDF) SetY(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetY(p.GetFloat("y", logic, 0.0))
	return pdf
}

// SetX maps json to gofpdf SetX function. Pass "x" as float object property in json logic.
// Default is "x": 0.0
func (p *JSONGOFPDF) SetX(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetX(p.GetFloat("x", logic, 0.0))
	return pdf
}

// SetXY maps json to gofpdf SetXY function. Pass "x" and "y" as float object properties in json logic.
// Defaults are "x": 0.0, "y": 0.0
func (p *JSONGOFPDF) SetXY(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetXY(p.GetFloat("x", logic, 0.0), p.GetFloat("y", logic, 0.0))
	return pdf
}

// Rect maps json to gofpdf Rect function. Pass "x", "y", "w", "h" float and "style" string object properties into json logic.
// Defaults are "x": 0.0, "y": 0.0, "w": 0.0, "h": 0.0, "style": ""
func (p *JSONGOFPDF) Rect(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.Rect(p.GetFloat("x", logic, 0.0), p.GetFloat("y", logic, 0.0), p.GetFloat("w", logic, 0.0), p.GetFloat("h", logic, 0.0), p.GetString("style", logic, ""))
	return pdf
}

// SetFooterFunc maps json to gofpdf SetFooterFunc function. Pass in an array of operation objects to have them be executed.
func (p *JSONGOFPDF) SetFooterFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf = p.RunArrayOperations(pdf, logic)
	})
	return pdf
}

// CellFormat maps json to gofpdf CellFormat function. Pass in "width" float, "height" float, "border" string, "text" string, "line" int, "align" string, "fill" boolean, "link" integer, "linkstr" string
// Defaults are "width": 0.0, "height": 0.0, "text": "", "border": "", "line": 0, "align": "L", "fill": false, "link": 0, "linkstr": ""
func (p *JSONGOFPDF) CellFormat(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	text := p.GetString("text", logic, "")
	text = strings.Replace(text, "<br>", "\n", -1)
	text = strings.Replace(text, "{nn}", cast.ToString(pdf.PageNo()), -1)
	if v := p.GetString("calculation", logic, ""); v != "" {
		text = p.Calculation(v, text)
	}

	text = p.Format(p.GetString("format", logic, ""), text)
	for index, value := range p.Globals {
		text = strings.Replace(text, index, cast.ToString(value), -1)
	}
	pdf.CellFormat(p.GetFloat("width", logic, 0.0), p.GetFloat("height", logic, 0.0), p.tr(text), p.GetString("border", logic, ""), p.GetInt("line", logic, 0), p.GetString("align", logic, "L"), p.GetBool("fill", logic, false), p.GetInt("link", logic, 0), p.GetString("linkstr", logic, ""))
	return pdf
}

// Cell maps json to gofpdf Cell function. Pass in "width" float, "height" float, "text" string object properties in json logic.
// Defaults are "width": 0.0, "height": 0.0, "text": ""
func (p *JSONGOFPDF) Cell(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.Cell(p.GetFloat("width", logic, 0.0), p.GetFloat("height", logic, 0.0), p.GetStringIndex("text", logic, ""))
	return pdf
}
