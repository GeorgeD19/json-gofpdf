package jsongofpdf

import (
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/GeorgeD19/json-logic-go"
	"github.com/GeorgeD19/securigroupgo/repository/data"
	"github.com/buger/jsonparser"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cast"
)

func New(options JSONGOFPDFOptions) (*JSONGOFPDF, error) {
	jsongofpdf := &JSONGOFPDF{}

	jsongofpdf.Logic = options.Logic
	jsongofpdf.Parser = options.Parser
	jsongofpdf.Form = options.Form
	jsongofpdf.Submission = options.Submission
	jsongofpdf.currentPage = 0
	jsongofpdf.DPI = 18

	return jsongofpdf, nil
}

// Apply is the entry function to parse logic and optional data
func (p *JSONGOFPDF) GetPDF() (opdf *gofpdf.Fpdf) {
	pdf := new(gofpdf.Fpdf)
	pdf = p.New(pdf, "{}")

	p.DocWidth, _ = pdf.GetPageSize()
	p.initY = pdf.GetY()

	// "" defaults to "cp1252" | This removes unwanted Â from special characters e.g. £
	p.tr = pdf.UnicodeTranslatorFromDescriptor("")

	result, _ := p.RunOperations(pdf, p.Logic, RowOptions{Index: 0})

	return result
}

// RunOperations will iterate through the array of operations and execute each
func (p *JSONGOFPDF) RunOperations(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	nRow = row
	jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			pdf, nRow = p.ParseObject(pdf, string(value), nRow)
			break
		}
	})

	return pdf, nRow
}

// ParseObject entry point
func (p *JSONGOFPDF) ParseObject(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	nRow = row
	jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		pdf, nRow = p.RunOperation(pdf, string(key), string(value), nRow)
		return nil
	})

	return pdf, nRow
}

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
		pdf = p.PreOperation(pdf, string(key), string(value))
		return nil
	})

	return pdf
}

func (p *JSONGOFPDF) PreOperation(pdf *gofpdf.Fpdf, name string, logic string) (opdf *gofpdf.Fpdf) {
	switch name {
	case "multicellformfield":
		p.PreMultiCellFormField(pdf, logic)
		break
	}
	return pdf
}

// RunOperation ensures that any operation ran doesn't crash the system if it doesn't exist
func (p *JSONGOFPDF) RunOperation(pdf *gofpdf.Fpdf, name string, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	p.CurrentX = pdf.GetX()
	p.CurrentY = pdf.GetY()

	// fmt.Printf("CurrentX: %v, CurrentY: %v, Operation: %v, Page: %v", p.CurrentX, p.CurrentY, name, p.currentPage)
	// fmt.Println("")

	nRow = row
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
	case "gety":
		pdf = p.GetY(pdf)
		break
	case "setinity":
		pdf, nRow = p.SetInitY(pdf, logic, nRow)
		break
	case "updatex":
		pdf = p.UpdateX(pdf, logic)
		break
	case "updatey":
		pdf = p.UpdateY(pdf, logic)
		break
	case "rowy":
		pdf = p.RowY(pdf, nRow)
		break
	case "setxy":
		pdf = p.SetXY(pdf, logic)
		break
	case "cell":
		pdf = p.Cell(pdf, logic, nRow)
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
		pdf, _ = p.SetHeaderFunc(pdf, logic, nRow)
		break
	case "setfooterfunc":
		pdf, _ = p.SetFooterFunc(pdf, logic)
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
	case "formfunc":
		pdf, _ = p.FormFunc(pdf, logic)
		break
	case "ln":
		pdf = p.Ln(pdf, logic)
		break
	case "image":
		pdf = p.Image(pdf, logic)
		break
	case "cellformfield":
		pdf = p.CellFormField(pdf, logic, nRow)
		break
	case "multicellformfield":
		pdf, nRow = p.MultiCellFormField(pdf, logic, nRow)
		break
	default:
		return pdf, nRow
	}
	return pdf, nRow
}

// ParseObjectValue entry point
func (p *JSONGOFPDF) ParseObjectValue(logic string, index RowOptions) (val []byte, dataType jsonparser.ValueType) {
	jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		val, dataType = p.RunValue(string(key), string(value), index)
		return nil
	})

	return val, dataType
}

func (p *JSONGOFPDF) RunValue(name string, logic string, row RowOptions) (val []byte, dataType jsonparser.ValueType) {
	switch name {
	// TODO Maybe worth using logic which opens up the entire logic library rather than just var and also means we don't have to merge the two libraries together
	case "logic":
		// uses json-logic and can read from parser data
		data := p.Parser.Parse()
		result, _ := jsonlogic.Apply(logic, data)
		switch v := result.(type) {
		case bool:
			return []byte(cast.ToString(v)), jsonparser.Boolean
		case int:
			return []byte(cast.ToString(v)), jsonparser.Number
		case float64:
			return []byte(cast.ToString(v)), jsonparser.Number
		case string:
			return []byte(cast.ToString(v)), jsonparser.String
		default:
			return []byte(cast.ToString(v)), jsonparser.String
		}
		break
	case "field":
		// Method to directly access a specific field value based on it's path or pass in a global variable within formfunc and it will return based on index
		result, _ := p.Field(logic, row)
		switch v := result.(type) {
		case bool:
			return []byte(cast.ToString(v)), jsonparser.Boolean
		case int:
			return []byte(cast.ToString(v)), jsonparser.Number
		case float64:
			return []byte(cast.ToString(v)), jsonparser.Number
		case string:
			return []byte(cast.ToString(v)), jsonparser.String
		default:
			return []byte(cast.ToString(v)), jsonparser.String
		}

		break
	case "form":
		result := p.GetFormValue(logic)
		return []byte(cast.ToString(result)), jsonparser.String
		break
	case "submission":
		result := p.GetSubmissionValue(logic)
		return []byte(cast.ToString(result)), jsonparser.String
		break
	default:
		return nil, jsonparser.NotExist
	}

	return nil, jsonparser.NotExist
}

// GetForm returns supported attributes from the passed form
func (p *JSONGOFPDF) GetFormValue(logic string) (res interface{}) {
	if p.Form != nil {
		switch logic {
		case "title":
			return p.Form.Title
			break
		case "created_by":
			user, err := data.GetUserItem(p.Form.CreatedBy)
			if err != nil {
				return ""
			}
			return cast.ToString(user.Name)
			break
		case "created_at":
			return cast.ToString(p.Form.CreatedAt.Format("02/01/2006 03:04:05"))
			break
		}
	}
	return nil
}

// GetForm returns supported attributes from the passed form
func (p *JSONGOFPDF) GetSubmissionValue(logic string) (res interface{}) {
	if p.Submission != nil {
		switch logic {
		case "created_by":
			user, err := data.GetUserItem(p.Submission.CreatedBy)
			if err != nil {
				return ""
			}
			return cast.ToString(user.Name)
			break
		case "created_at":
			return cast.ToString(p.Submission.CreatedAt.Format("02/01/2006 03:04:05"))
			break
		}
	}
	return nil
}

func (p *JSONGOFPDF) Field(logic string, row RowOptions) (res interface{}, err error) {
	result := logic

	for k := range p.Parser.FieldRegistry {
		if k == p.RowIndex {
			item := p.Parser.FieldRegistry[p.RowIndex]

			result = strings.Replace(result, "{field:title}", strings.Replace(cast.ToString(item.Title), "<br>", "\n", -1), -1)
			result = strings.Replace(result, "{field:value}", cast.ToString(item.Value), -1)
		}
	}

	return result, err
}

func (p *JSONGOFPDF) CellFormField(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf) {
	attribute := p.GetString("attribute", logic, "")
	target := p.GetString("target", logic, "")

	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0)
	border := p.GetString("border", logic, "")
	line := p.GetInt("line", logic, 0)
	align := p.GetString("align", logic, "L")
	fill := p.GetBool("fill", logic, false)
	link := p.GetInt("link", logic, 0)
	linkStr := p.GetString("linkstr", logic, "")

	field := p.Parser.FieldRegistry[row.Index]
	if target != "" {
		for _, v := range p.Parser.FieldRegistry {
			if v.PathString == target || v.Key == target {
				field = v
			}
		}
	}

	CurrentX := pdf.GetX()

	switch attribute {
	case "title":
		pdf.CellFormat(width, height, p.tr(strings.Replace(field.Title, "<br>", "\n", -1)), border, line, align, fill, link, linkStr)
		break
	case "value":
		pdf.CellFormat(width, height, p.tr(strings.Replace(cast.ToString(field.Value), "<br>", "\n", -1)), border, line, align, fill, link, linkStr)

		// fmt.Println("Media displayed below")
		// fmt.Println(field.Media)
		for _, image := range field.Media {

			// For any media against the field
			imageDecoded, _ := hex.DecodeString(image.Data)
			options := gofpdf.ImageOptions{
				ReadDpi:   false,
				ImageType: image.Type,
			}

			// Pass binary image into PDF
			pdf.RegisterImageOptionsReader(string(p.MediaIndex), options, bytes.NewReader(imageDecoded))
			pdf.ImageOptions(string(p.MediaIndex), CurrentX, pdf.GetY(), width, 0, true, options, -1, "")
			p.MediaIndex++
		}
	}

	return pdf
}

func (p *JSONGOFPDF) MultiCellFormField(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	nRow = row

	attribute := p.GetString("attribute", logic, "")
	target := p.GetString("target", logic, "")

	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0) // Line height of each cell, not cell height
	border := p.GetString("border", logic, "")
	align := p.GetString("align", logic, "L")
	fill := p.GetBool("fill", logic, false)

	field := p.Parser.FieldRegistry[row.Index]
	if target != "" {
		for _, v := range p.Parser.FieldRegistry {
			if v.PathString == target || v.Key == target {
				field = v
			}
		}
	}

	if p.RowHeight > height {
		height = p.RowHeight
	}

	if field.Key != "" {
		switch attribute {
		case "title":
			text := p.tr(strings.Replace(field.Title, "<br>", "\n", -1))
			cellList := pdf.SplitLines([]byte(text), width)
			cellCount := float64(len(cellList))

			if cellCount > p.RowCells {
				p.RowCells = cellCount
			}

			cellHeight := height / cellCount

			// Line height = number of cells / total height
			pdf.MultiCell(width, cellHeight, p.tr(strings.Replace(field.Title, "<br>", "\n", -1)), border, align, fill)
			break
		case "value":

			CurrentX := pdf.GetX()

			text := p.tr(strings.Replace(cast.ToString(field.Value), "<br>", "\n", -1))
			cellList := pdf.SplitLines([]byte(text), width)
			cellCount := float64(len(cellList))
			if cellCount < 1 {
				cellCount = 1
			}

			if cellCount > p.RowCells {
				p.RowCells = cellCount
			}

			cellHeight := height / cellCount

			pdf.MultiCell(width, cellHeight, text, border, align, fill)

			for _, image := range field.Media {

				// For any media against the field
				imageDecoded, _ := hex.DecodeString(image.Data)
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
				pdf.ImageOptions(string(p.MediaIndex), CurrentX, pdf.GetY(), imageWidth, imageHeight, true, options, -1, "")
				p.MediaIndex++
			}
		}
	}

	// postRenderHeight := pdf.GetY()

	// cellHeight := postRenderHeight - preRenderHeight
	// p.PrevCellHeight = cellHeight

	CurrentY := pdf.GetY()
	if CurrentY > p.NextY {
		p.NextY = CurrentY
	}

	p.NewPage = false

	return pdf, nRow
}

func (p *JSONGOFPDF) UpdateX(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	width := p.GetFloat("width", logic, 0.0)
	pdf.SetX(pdf.GetX() + width)
	return pdf
}

func (p *JSONGOFPDF) UpdateY(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	height := p.GetFloat("height", logic, 0.0)
	pdf.SetY(pdf.GetY() + height)
	return pdf
}

func (p *JSONGOFPDF) FormFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	main := p.GetString("main", logic, "")
	alternative := p.GetString("alternative", logic, "")

	n := 2
	if main != "" {
		if alternative != "" {
			for x := 0; x < len(p.Parser.FieldRegistry); x++ {

				nRow.PrevCellHeight = 0
				nRow.Index = x
				p.RowIndex = x
				p.PrevCellHeight = 0
				p.RowHeight = 0
				p.RowCells = 0.0

				p.PreOperations(pdf, main)

				if x%n == 0 {
					pdf, nRow = p.RunOperations(pdf, alternative, nRow)
				} else {
					pdf, nRow = p.RunOperations(pdf, main, nRow)
				}
				if nRow.NextY > nRow.CurrentY {
					nRow.CurrentY = nRow.NextY
				}
				p.CurrentRowY = pdf.GetY()
			}
		} else {
			for x := 0; x < len(p.Parser.FieldRegistry); x++ {

				nRow.PrevCellHeight = 0
				p.PrevCellHeight = 0
				nRow.Index = x
				p.RowIndex = x
				p.RowHeight = 0
				p.RowCells = 0.0

				p.PreOperations(pdf, main)

				pdf, nRow = p.RunOperations(pdf, main, nRow)
				p.CurrentRowY = pdf.GetY()
			}
		}
	}

	return pdf, nRow
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

	image, _ := p.GetImage(src)
	imageContent := strings.TrimPrefix(image.Data, "0x")
	imageDecoded, _ := hex.DecodeString(imageContent)

	options := gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: image.Type,
	}

	// Pass binary image into PDF
	pdf.RegisterImageOptionsReader(name, options, bytes.NewReader(imageDecoded))
	pdf.ImageOptions(name, x, y, width, height, flow, options, link, linkStr)
	return pdf
}

func (p *JSONGOFPDF) SetHeaderFunc(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	pdf.SetHeaderFunc(func() {
		nRow = row
		p.NewPage = true
		p.currentPage++
		pdf, nRow = p.RunOperations(pdf, logic, nRow)
		p.CurrentRowY = pdf.GetY()
		p.CurrentY = pdf.GetY()
		p.HeaderHeight = pdf.GetY()
	})

	return pdf, nRow
}

func (p *JSONGOFPDF) GetStringIndex(name string, logic string, fallback string, row RowOptions) (value string) {
	result := fallback
	attribute, _, _, err := p.GetAttributeIndex(name, logic, true, row)
	if err == nil {
		result = cast.ToString(attribute)
	}
	return result
}

func (p *JSONGOFPDF) GetAttribute(name string, logic string, debug bool) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	return p.GetAttributeIndex(name, logic, debug, RowOptions{Index: 0})
}

func (p *JSONGOFPDF) GetAttributeIndex(name string, logic string, debug bool, row RowOptions) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	value, dataType, offset, err = jsonparser.Get([]byte(logic), name)
	if dataType == jsonparser.Object {
		value, dataType = p.ParseObjectValue(string(value), row)
	}
	return value, dataType, offset, err
}
