package jsongofpdf

import (
	"fmt"
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
// a4 (w:210.00155555555557, h:841.89) // docWidth - 30 (margin from left / right) = 180.00155555555557
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

// New passes the orientation, unit, size and dir object properties to the gofpdf New function creating a new pdf
func (p *JSONGOFPDF) New(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	orientation := p.GetString("orientation", logic, "P")
	unit := p.GetString("unit", logic, "mm")
	size := p.GetString("size", logic, "A4")
	directory := p.GetString("dir", logic, "")
	return gofpdf.New(orientation, unit, size, directory)
}

func (p *JSONGOFPDF) SetInitY(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	nRow = row
	pdf.SetY(p.initY)
	nRow.NextY = p.initY
	return pdf, nRow
}

// SetMargins passes the left, top and right object properties to the gofpdf SetMargins function affecting the working pdf
func (p *JSONGOFPDF) SetMargins(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	left := p.GetFloat("left", logic, 0.0)
	top := p.GetFloat("top", logic, 0.0)
	right := p.GetFloat("right", logic, 0.0)
	pdf.SetMargins(left, top, right)
	return pdf
}

func (p *JSONGOFPDF) SetCellMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	margin := p.GetFloat("margin", logic, 0.0)
	pdf.SetCellMargin(margin)
	return pdf
}

func (p *JSONGOFPDF) SetTopMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	margin := p.GetFloat("margin", logic, 0.0)
	pdf.SetTopMargin(margin)
	return pdf
}

func (p *JSONGOFPDF) SetRightMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	margin := p.GetFloat("margin", logic, 0.0)
	pdf.SetRightMargin(margin)
	return pdf
}

func (p *JSONGOFPDF) SetAutoPageBreak(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	auto := p.GetBool("auto", logic, true)
	margin := p.GetFloat("margin", logic, 15.0)
	pdf.SetAutoPageBreak(auto, margin)
	return pdf
}

// SetDisplayMode passes the zoom and layout object properties to the gofpdf SetDisplayMode function affecting the working pdf
func (p *JSONGOFPDF) SetDisplayMode(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	zoomStr := p.GetString("zoom", logic, "")
	layoutStr := p.GetString("layout", logic, "")
	pdf.SetDisplayMode(zoomStr, layoutStr)
	return pdf
}

// SetDefaultCompression passes the compress object property to the gofpdf SetDefaultCompression function affecting the working pdf
// func (p *JSONGOFPDF) SetDefaultCompression(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
// 	compress := p.GetBool("compress", logic, false)
// 	pdf.SetDefaultCompression(compress)
// 	return pdf
// }

func (p *JSONGOFPDF) Ln(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	height := p.GetFloat("height", logic, -1.0)
	pdf.Ln(height)
	return pdf
}

func (p *JSONGOFPDF) SetDrawColor(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	r := p.GetInt("r", logic, 0)
	g := p.GetInt("g", logic, 0)
	b := p.GetInt("b", logic, 0)
	pdf.SetDrawColor(r, g, b)
	return pdf
}

func (p *JSONGOFPDF) SetFillColor(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	r := p.GetInt("r", logic, 0)
	g := p.GetInt("g", logic, 0)
	b := p.GetInt("b", logic, 0)
	pdf.SetFillColor(r, g, b)
	return pdf
}

func (p *JSONGOFPDF) SetTextColor(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	r := p.GetInt("r", logic, 0)
	g := p.GetInt("g", logic, 0)
	b := p.GetInt("b", logic, 0)
	pdf.SetTextColor(r, g, b)
	return pdf
}

// AddPage adds a new page to the pdf
func (p *JSONGOFPDF) AddPage(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.AddPage()
	return pdf
}

// SetFont sets the font for the pdf
func (p *JSONGOFPDF) SetFont(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	family := p.GetString("family", logic, "Arial")
	style := p.GetString("style", logic, "")
	size := p.GetFloat("size", logic, 8.0)
	pdf.SetFont(family, style, size)
	return pdf
}

func (p *JSONGOFPDF) AliasNbPages(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	aliasStr := p.GetString("alias", logic, "")
	pdf.AliasNbPages(aliasStr)
	return pdf
}

func (p *JSONGOFPDF) SetLeftMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	margin := p.GetFloat("margin", logic, 0.0)
	pdf.SetLeftMargin(margin)
	return pdf
}

func (p *JSONGOFPDF) GetY(pdf *gofpdf.Fpdf) (opdf *gofpdf.Fpdf) {
	fmt.Println(pdf.GetY())
	return pdf
}

func (p *JSONGOFPDF) SetY(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	y := p.GetFloat("y", logic, 0.0)
	pdf.SetY(y)
	return pdf
}

func (p *JSONGOFPDF) SetX(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	x := p.GetFloat("x", logic, 0.0)
	pdf.SetX(x)
	return pdf
}

func (p *JSONGOFPDF) SetXY(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	x := p.GetFloat("x", logic, 0.0)
	y := p.GetFloat("y", logic, 0.0)
	pdf.SetXY(x, y)
	return pdf
}

func (p *JSONGOFPDF) CellFormat(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0)
	text := p.GetString("text", logic, "")
	text = strings.Replace(text, "{nn}", cast.ToString(pdf.PageNo()), -1)
	border := p.GetString("border", logic, "")
	line := p.GetInt("line", logic, 0)
	align := p.GetString("align", logic, "L")
	fill := p.GetBool("fill", logic, false)
	link := p.GetInt("link", logic, 0)
	linkStr := p.GetString("linkstr", logic, "")
	pdf.CellFormat(width, height, text, border, line, align, fill, link, linkStr)
	return pdf
}

func (p *JSONGOFPDF) SetFooterFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	pdf.SetFooterFunc(func() {
		pdf, nRow = p.RunOperations(pdf, logic, RowOptions{Index: 0})
	})
	return pdf, nRow
}

func (p *JSONGOFPDF) Cell(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf) {
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0)
	text := p.GetStringIndex("text", logic, "", row)
	pdf.Cell(width, height, text)
	return pdf
}

// RowY sets pdf Y to RowY position
func (p *JSONGOFPDF) RowY(pdf *gofpdf.Fpdf, row RowOptions) (opdf *gofpdf.Fpdf) {
	pdf.SetY(p.CurrentRowY)
	return pdf
}
