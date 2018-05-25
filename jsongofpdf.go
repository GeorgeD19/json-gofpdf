package jsongofpdf

import (
	"github.com/buger/jsonparser"
	"github.com/jung-kurt/gofpdf"
)

// New creates a new jsongofpdf instance.
func New(options JSONGOFPDFOptions) (*JSONGOFPDF, error) {
	jsongofpdf := &JSONGOFPDF{}

	jsongofpdf.Logic = options.Logic
	jsongofpdf.Tables = options.Tables
	jsongofpdf.Globals = options.Globals

	jsongofpdf.DPI = 18

	return jsongofpdf, nil
}

// GetPDF parses logic and generates a pdf.
func (p *JSONGOFPDF) GetPDF() (opdf *gofpdf.Fpdf) {
	pdf := new(gofpdf.Fpdf)
	pdf = p.New(pdf, "{}")

	p.DocWidth, _ = pdf.GetPageSize()
	p.initY = pdf.GetY()

	// "" defaults to "cp1252" | This removes unwanted Â from special characters e.g. £
	p.tr = pdf.UnicodeTranslatorFromDescriptor("")

	result := p.RunArrayOperations(pdf, p.Logic)

	return result
}

// RunArrayOperations will iterate through the array of operations and execute each one.
func (p *JSONGOFPDF) RunArrayOperations(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			pdf = p.RunObjectOperations(pdf, string(value))
			break
		}
	})
	return pdf
}

// RunObjectOperations entry point
func (p *JSONGOFPDF) RunObjectOperations(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		pdf = p.RunOperation(pdf, string(key), string(value))
		return nil
	})
	return pdf
}

// RunOperation ensures that any operation ran doesn't crash the system if it doesn't exist
func (p *JSONGOFPDF) RunOperation(pdf *gofpdf.Fpdf, name string, logic string) (opdf *gofpdf.Fpdf) {
	p.CurrentY = pdf.GetY()

	switch name {
	case "new":
		pdf = p.New(pdf, logic)
		break
	case "addpage":
		pdf = p.AddPage(pdf, logic)
		break
	case "setfont":
		pdf = p.SetFont(pdf, logic)
		break
	case "setx":
		pdf = p.SetX(pdf, logic)
		break
	case "sety":
		pdf = p.SetY(pdf, logic)
		break
	case "setinity":
		pdf = p.SetInitY(pdf, logic)
		break
	case "updatex":
		pdf = p.UpdateX(pdf, logic)
		break
	case "updatey":
		pdf = p.UpdateY(pdf, logic)
		break
	case "rowy":
		pdf = p.RowY(pdf)
		break
	case "setxy":
		pdf = p.SetXY(pdf, logic)
		break
	case "cell":
		pdf = p.Cell(pdf, logic)
		break
	case "cellformat":
		pdf = p.CellFormat(pdf, logic)
		break
	case "setmargins":
		pdf = p.SetMargins(pdf, logic)
		break
	case "setautopagebreak":
		pdf = p.SetAutoPageBreak(pdf, logic)
		break
	case "aliasnbpages":
		pdf = p.AliasNbPages(pdf, logic)
		break
	case "setheaderfunc":
		pdf = p.SetHeaderFunc(pdf, logic)
		break
	case "setfooterfunc":
		pdf = p.SetFooterFunc(pdf, logic)
		break
	case "settopmargin":
		pdf = p.SetTopMargin(pdf, logic)
		break
	case "setleftmargin":
		pdf = p.SetLeftMargin(pdf, logic)
		break
	case "settextcolor":
		pdf = p.SetTextColor(pdf, logic)
		break
	case "setfillcolor":
		pdf = p.SetFillColor(pdf, logic)
		break
	case "setdrawcolor":
		pdf = p.SetDrawColor(pdf, logic)
		break
	case "tablefunc":
		pdf = p.TableFunc(pdf, logic)
		break
	case "ln":
		pdf = p.Ln(pdf, logic)
		break
	case "image":
		pdf = p.Image(pdf, logic)
		break
	case "multicell":
		pdf = p.MultiCell(pdf, logic)
		break
	default:
		return pdf
	}
	return pdf
}
