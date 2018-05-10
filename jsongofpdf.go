package jsongofpdf

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"strings"

	"github.com/GeorgeD19/json-logic-go"
	"github.com/buger/jsonparser"
	"github.com/h2non/filetype"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cast"
)

func New(options JSONGOFPDFOptions) (*JSONGOFPDF, error) {
	jsongofpdf := &JSONGOFPDF{}

	jsongofpdf.Logic = options.Logic
	jsongofpdf.Parser = options.Parser
	jsongofpdf.Form = options.Form
	jsongofpdf.currentPage = 0

	return jsongofpdf, nil
}

// Apply is the entry function to parse logic and optional data
func (p *JSONGOFPDF) GetPDF() (opdf *gofpdf.Fpdf) {
	pdf := new(gofpdf.Fpdf)
	pdf = p.New(pdf, "{}")

	p.DocWidth, _ = pdf.GetPageSize()
	p.initY = pdf.GetY()

	// fmt.Println(p.DocWidth)

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

func (p *JSONGOFPDF) New(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	orientation := p.GetString("orientation", logic, "P")
	unit := p.GetString("unit", logic, "mm")
	size := p.GetString("size", logic, "A4")
	directory := p.GetString("dir", logic, "")
	return gofpdf.New(orientation, unit, size, directory)
}

// RunOperation ensures that any operation ran doesn't crash the system if it doesn't exist
func (p *JSONGOFPDF) RunOperation(pdf *gofpdf.Fpdf, name string, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	// fmt.Println(name)
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

// func (p *JSONGOFPDF) ParseObjectValue(logic string, index int) (val []byte, dataType jsonparser.ValueType) {
// initX = pdf.GetX()
// 		initY = pdf.GetY() + marginRw
// 		pdf.SetY(initY)
// }

func (p *JSONGOFPDF) SetInitY(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	nRow = row
	pdf.SetY(p.initY)
	nRow.NextY = p.initY
	return pdf, nRow
}

// ParseObject entry point
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
		result := p.GetForm(logic)
		return []byte(cast.ToString(result)), jsonparser.String
		break
	default:
		return nil, jsonparser.NotExist
	}

	return nil, jsonparser.NotExist
}

func (p *JSONGOFPDF) GetForm(logic string) (res interface{}) {

	if p.Form != nil {
		switch logic {
		case "title":
			return p.Form.Title
			break
		case "created_by":
			return cast.ToString(p.Form.CreatedBy)
			break
		case "created_at":
			return cast.ToString(p.Form.CreatedAt)
			break
		}
	}

	return nil
}

func (p *JSONGOFPDF) Field(logic string, row RowOptions) (res interface{}, err error) {
	// fmt.Println(index)
	// fmt.Println(logic)

	result := logic

	for k := range p.Parser.FieldRegistry {
		if k == row.Index {
			item := p.Parser.FieldRegistry[row.Index]

			result = strings.Replace(result, "{field:title}", strings.Replace(cast.ToString(item.Title), "<br>", "\n", -1), -1)
			result = strings.Replace(result, "{field:value}", cast.ToString(item.Value), -1)
		}
	}

	// rowValue, rowIsset := p.Parser.FieldRegistry[index]
	// if rowIsset {
	// 	// Try replace {title} with rowValue.title
	// 	// rowValue.Title
	// }

	// targetValue, targetIsset := p.Parser.FieldRegistry[logic]
	// if targetIsset {

	// }

	// We get the current looped field

	// Replace the global field

	// We get the field requested if it is not a global field

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

	switch attribute {
	case "title":
		pdf.CellFormat(width, height, p.tr(strings.Replace(field.Title, "<br>", "\n", -1)), border, line, align, fill, link, linkStr)

		break
	case "value":
		pdf.CellFormat(width, height, p.tr(cast.ToString(field.Value)), border, line, align, fill, link, linkStr)

		// for _, image := range field.Media {

		// 	// For any media against the field
		// 	ImageDecoded, _ := hex.DecodeString(image.Data)

		// }

	}
	return pdf
}

func (p *JSONGOFPDF) MultiCellFormField(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	nRow = row

	// if p.NewPage {
	// 	pdftest := pdf.GetY()
	// 	fmt.Println(pdftest)
	// 	if p.currentPage > 1 {
	// 		// pdf.SetY(nRow.NextY)
	// 	}
	// }

	attribute := p.GetString("attribute", logic, "")
	target := p.GetString("target", logic, "")

	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0)
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

	if nRow.PrevCellHeight > height {
		height = nRow.PrevCellHeight
	}

	pdftest := pdf.GetY()
	if pdftest == 0 {
		// So it seems on creating a new page we lose our margin, and x/y positioning
		pdf.SetXY(width+15, p.HeaderHeight)

	}

	if field.Key != "" {
		switch attribute {
		case "title":
			pdf.MultiCell(width, height, p.tr(strings.Replace(field.Title, "<br>", "\n", -1)), border, align, fill)
			break
		case "value":
			pdf.MultiCell(width, height, p.tr(cast.ToString(field.Value)), border, align, fill)
			// for _, image := range field.Media {
			// 	// For any media against the field
			// 	ImageDecoded, _ := hex.DecodeString(image.Data)
			// }
		}
	}

	if p.NewPage {
		nRow.CurrentY = p.initY
		nRow.NextY = p.initY
	}

	pdfY := pdf.GetY()

	if pdfY >= nRow.NextY {
		nRow.NextY = pdf.GetY()
	}

	// We need to track the height of the previous cell in the p object
	if !p.NewPage {
		nRow.PrevCellHeight = pdf.GetY() - nRow.CurrentY
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

func (p *JSONGOFPDF) RowY(pdf *gofpdf.Fpdf, row RowOptions) (opdf *gofpdf.Fpdf) {
	pdf.SetY(row.CurrentY)
	return pdf
}

func (p *JSONGOFPDF) Cell(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf) {
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0)
	text := p.GetStringIndex("text", logic, "", row)
	pdf.Cell(width, height, text)
	return pdf
}

func (p *JSONGOFPDF) FormFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	main := p.GetString("main", logic, "")
	alternative := p.GetString("alternative", logic, "")
	nRow.Index = 0

	n := 2
	if main != "" {
		// FormFunc needs it's own index counter for form fields so it doesn't interfere on any other page / section
		if alternative != "" {
			nRow.CurrentY = pdf.GetY()
			// nRow := RowOptions{Index: 0, CurrentY: pdf.GetY()}
			for x := 0; x < len(p.Parser.FieldRegistry); x++ {
				nRow.PrevCellHeight = 0
				nRow.Index = x
				if x%n == 0 {
					pdf, nRow = p.RunOperations(pdf, main, nRow)
				} else {
					pdf, nRow = p.RunOperations(pdf, alternative, nRow)
				}
				if nRow.NextY > nRow.CurrentY {
					nRow.CurrentY = nRow.NextY
				}

			}
		} else {
			nRow.CurrentY = pdf.GetY()
			// nRow := RowOptions{Index: 0, CurrentY: pdf.GetY()}
			for x := 0; x < len(p.Parser.FieldRegistry); x++ {
				nRow.PrevCellHeight = 0
				nRow.Index = x
				pdf, nRow = p.RunOperations(pdf, main, nRow)
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

func (p *JSONGOFPDF) SetAutoPageBreak(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	auto := p.GetBool("auto", logic, true)
	margin := p.GetFloat("margin", logic, 15.0)
	pdf.SetAutoPageBreak(auto, margin)
	return pdf
}

func (p *JSONGOFPDF) AliasNbPages(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	aliasStr := p.GetString("alias", logic, "")
	pdf.AliasNbPages(aliasStr)
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

func (p *JSONGOFPDF) SetTopMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	margin := p.GetFloat("margin", logic, 0.0)
	pdf.SetTopMargin(margin)
	return pdf
}

func (p *JSONGOFPDF) SetLeftMargin(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	margin := p.GetFloat("margin", logic, 0.0)
	pdf.SetLeftMargin(margin)
	return pdf
}

func (p *JSONGOFPDF) SetHeaderFunc(pdf *gofpdf.Fpdf, logic string, row RowOptions) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	pdf.SetHeaderFunc(func() {
		nRow = row
		p.NewPage = true
		p.currentPage++
		// fmt.Println(p.NewPage)
		// fmt.Println(p.currentPage)
		// pdf.SetY(p.initY)
		// nRow.CurrentY = p.initY
		pdf, nRow = p.RunOperations(pdf, logic, nRow)
		p.HeaderHeight = pdf.GetY()
		// pdf.SetY(pdf.GetY())
		// fmt.Println(pdf.GetY())
	})

	return pdf, nRow
}

func (p *JSONGOFPDF) SetFooterFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, nRow RowOptions) {
	pdf.SetFooterFunc(func() {
		pdf, nRow = p.RunOperations(pdf, logic, RowOptions{Index: 0})
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

func (p *JSONGOFPDF) GetBool(name string, logic string, fallback bool) (value bool) {
	result := fallback
	attribute, _, _, err := p.GetAttribute(name, logic, false)
	if err == nil {
		result, _ = jsonparser.ParseBoolean(attribute)
	}
	return result
}

func (p *JSONGOFPDF) GetFloat(name string, logic string, fallback float64) (value float64) {
	result := fallback
	attribute, _, _, err := p.GetAttribute(name, logic, false)
	if err == nil {
		result = cast.ToFloat64(cast.ToString(attribute))
	}
	return result
}

func (p *JSONGOFPDF) GetInt(name string, logic string, fallback int) (value int) {
	result := fallback
	attribute, _, _, err := p.GetAttribute(name, logic, false)

	if err == nil {
		result = cast.ToInt(cast.ToString(attribute))
	}
	return result
}

func (p *JSONGOFPDF) GetString(name string, logic string, fallback string) (value string) {
	return p.GetStringIndex(name, logic, fallback, RowOptions{Index: 0})
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

	// TODO Parse variables e.g. we can get a variable but it may be something like {"var": "something"} and it should return the correct value
	value, dataType, offset, err = jsonparser.Get([]byte(logic), name)
	if dataType == jsonparser.Object {

		// p.RunOperation(pdf, string(key), string(value))
		// We need to run a middle operator that only returns values

		value, dataType = p.ParseObjectValue(string(value), row)

		// value, dataType = p.ParseAttribute(string(value))
	}

	// TODO Determine if the getvalues has actually given something or if we should take the straight out value? or maybe that's already done as getvalues detects the type
	// attr := jsonlogic.GetValues(string(value), `{"test":"4"}`)
	// value = []byte(attr)
	// fmt.Println(cast.ToString(value))
	// fmt.Println(cast.ToInt(cast.ToString(attr)))

	// Here we can get the attribute from the object and if it matches any logic operations then we run them
	return value, dataType, offset, err
}

type ImageFile struct {
	Data   string
	Width  int
	Height int
	Type   string
	Mime   string
}

// GetImage returns a File type containing the hex version of the image with associated meta information
func (p *JSONGOFPDF) GetImage(FileName string) (f ImageFile, Err error) {

	FoundFile := ImageFile{}

	// Get file contents
	FileData, _ := ioutil.ReadFile(FileName)

	FileMeta, err := filetype.Match(FileData)
	if err != nil {
		return FoundFile, err
	}

	// Compatible with MSSQL binary storage
	FileContent := hex.EncodeToString(FileData)
	FileContent = "0x" + FileContent

	FoundFile.Data = FileContent
	FoundFile.Type = FileMeta.Extension
	FoundFile.Mime = FileMeta.MIME.Value

	File, err := os.Open(FileName)
	defer File.Close()
	if err != nil {
		return FoundFile, err
	}

	head := make([]byte, 261)
	File.Read(head)

	// TODO Define err message
	if filetype.IsImage(head) == false {
		return FoundFile, err
	}

	// Only parse for supported functions
	switch FoundFile.Type {
	case "jpg":
		FoundFile.Width, FoundFile.Height = p.GetJpgDimensions(File)
	case "gif":
		FoundFile.Width, FoundFile.Height = p.GetGifDimensions(File)
	case "png":
		FoundFile.Width, FoundFile.Height = p.GetPngDimensions(File)
	case "bmp":
		FoundFile.Width, FoundFile.Height = p.GetBmpDimensions(File)
	default:
		return FoundFile, err
	}

	return FoundFile, nil
}

func (p *JSONGOFPDF) GetJpgDimensions(file *os.File) (width int, height int) {
	fi, _ := file.Stat()
	fileSize := fi.Size()

	position := int64(4)
	bytes := make([]byte, 4)
	file.ReadAt(bytes[:2], position)
	length := int(bytes[0]<<8) + int(bytes[1])
	for position < fileSize {
		position += int64(length)
		file.ReadAt(bytes, position)
		length = int(bytes[2])<<8 + int(bytes[3])
		if (bytes[1] == 0xC0 || bytes[1] == 0xC2) && bytes[0] == 0xFF && length > 7 {
			file.ReadAt(bytes, position+5)
			width = int(bytes[2])<<8 + int(bytes[3])
			height = int(bytes[0])<<8 + int(bytes[1])
			return
		}
		position += 2
	}
	return 0, 0
}

func (p *JSONGOFPDF) GetGifDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 4)
	file.ReadAt(bytes, 6)
	width = int(bytes[0]) + int(bytes[1])*256
	height = int(bytes[2]) + int(bytes[3])*256
	return
}

func (p *JSONGOFPDF) GetBmpDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 18)
	width = int(bytes[3])<<24 | int(bytes[2])<<16 | int(bytes[1])<<8 | int(bytes[0])
	height = int(bytes[7])<<24 | int(bytes[6])<<16 | int(bytes[5])<<8 | int(bytes[4])
	return
}

func (p *JSONGOFPDF) GetPngDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 16)
	width = int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	height = int(bytes[4])<<24 | int(bytes[5])<<16 | int(bytes[6])<<8 | int(bytes[7])
	return
}
