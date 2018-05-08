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

	return jsongofpdf, nil
}

// Apply is the entry function to parse logic and optional data
func (p *JSONGOFPDF) GetPDF() (opdf *gofpdf.Fpdf, errs error) {
	pdf := new(gofpdf.Fpdf)
	result, err := p.RunOperations(pdf, p.Logic)

	return result, err
}

// RunOperations will iterate through the array of operations and execute each
func (p *JSONGOFPDF) RunOperations(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, err error) {

	parseErr := err

	jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			pdf, err = p.ParseObject(pdf, string(value))
			if err != nil {
				parseErr = err
			}
			break
		}

	})

	if parseErr != nil {
		return nil, parseErr
	}

	return pdf, nil
}

// RunOperation ensures that any operation ran doesn't crash the system if it doesn't exist
func (p *JSONGOFPDF) RunOperation(pdf *gofpdf.Fpdf, name string, logic string) (opdf *gofpdf.Fpdf, err error) {
	// fmt.Println(name)
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
	case "setxy":
		pdf = p.SetXY(pdf, logic)
		break
	case "cell":
		pdf = p.Cell(pdf, logic)
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
		pdf, _ = p.SetHeaderFunc(pdf, logic)
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
	case "bodyfunc":
		pdf, _ = p.BodyFunc(pdf, logic)
		break
	case "ln":
		pdf = p.Ln(pdf, logic)
		break
	case "image":
		pdf = p.Image(pdf, logic)
		break
	case "var":
		pdf = p.Var(pdf, logic)
		break
	default:
		return pdf, ErrInvalidOperation
	}
	return pdf, nil
}

func (p *JSONGOFPDF) New(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	orientation := p.GetString("orientation", logic, "P")
	unit := p.GetString("unit", logic, "mm")
	size := p.GetString("size", logic, "A4")
	directory := p.GetString("dir", logic, "")
	return gofpdf.New(orientation, unit, size, directory)
}

func (p *JSONGOFPDF) Var(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	// fmt.Println("Var")
	// fmt.Println(logic) // jsonlogic.Apply()

	// docWidth, _ := pdf.GetPageSize()
	// docWidth = docWidth - (marginH * 2)
	// tr := pdf.UnicodeTranslatorFromDescriptor("")

	// fmt.Println(data)
	// // Try extract the variable from the data
	// for x := 0; x < len(data); x++ {
	// 	if data[x].Key == logic {
	// 		fmt.Println(data[x].Value)
	// 		// 		pdf.MultiCell(docWidth/2, lineHt, tr(data[x].Field), "", "L", true)
	// 	}
	// }

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

func (p *JSONGOFPDF) BodyFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, err error) {
	for x := 0; x < len(p.Parser.FieldRegistry); x++ {
		pdf, err = p.RunOperations(pdf, logic)
	}
	return pdf, err
}

func (p *JSONGOFPDF) SetFillColor(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	r := p.GetInt("r", logic, 0)
	g := p.GetInt("r", logic, 0)
	b := p.GetInt("r", logic, 0)
	pdf.SetFillColor(r, g, b)
	return pdf
}

func (p *JSONGOFPDF) SetTextColor(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	r := p.GetInt("r", logic, 0)
	g := p.GetInt("r", logic, 0)
	b := p.GetInt("r", logic, 0)
	pdf.SetTextColor(r, g, b)
	return pdf
}

func (p *JSONGOFPDF) AddPage(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	pdf.AddPage()
	return pdf
}

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

func (p *JSONGOFPDF) Cell(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf) {
	width := p.GetFloat("width", logic, 0.0)
	height := p.GetFloat("height", logic, 0.0)
	text := p.GetString("text", logic, "")
	pdf.Cell(width, height, text)
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

func (p *JSONGOFPDF) SetHeaderFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, err error) {
	pdf.SetHeaderFunc(func() {
		pdf, err = p.RunOperations(pdf, logic)
	})
	return pdf, err
}

func (p *JSONGOFPDF) SetFooterFunc(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, err error) {
	pdf.SetFooterFunc(func() {
		pdf, err = p.RunOperations(pdf, logic)
	})
	return pdf, err
}

// ParseObject entry point
func (p *JSONGOFPDF) ParseObject(pdf *gofpdf.Fpdf, logic string) (opdf *gofpdf.Fpdf, err error) {

	err = jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		pdf, err = p.RunOperation(pdf, string(key), string(value))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return pdf, err
	}

	return pdf, nil
}

func (p *JSONGOFPDF) GetBool(name string, logic string, fallback bool) (value bool) {
	result := fallback
	attribute, _, _, err := p.GetAttribute(name, logic, false)
	if err == nil {
		result = cast.ToBool(attribute)
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
	result := fallback
	attribute, _, _, err := p.GetAttribute(name, logic, true)
	if err == nil {
		result = cast.ToString(attribute)
	}
	return result
}

func (p *JSONGOFPDF) ParseAttribute(logic string) (value []byte, dataType jsonparser.ValueType) {
	// Here we can use json-logic-go to parse the attribute and return the value as interface{}
	result, _ := jsonlogic.Run(logic)

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

}

func (p *JSONGOFPDF) GetAttribute(name string, logic string, debug bool) (value []byte, dataType jsonparser.ValueType, offset int, err error) {

	// TODO Parse variables e.g. we can get a variable but it may be something like {"var": "something"} and it should return the correct value
	value, dataType, offset, err = jsonparser.Get([]byte(logic), name)
	if dataType == jsonparser.Object {
		// p.RunOperation(pdf, string(key), string(value))
		value, dataType = p.ParseAttribute(string(value))
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
