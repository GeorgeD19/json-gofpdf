package jsongofpdf

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/GeorgeD19/libraries/alpaca"
	jsonlogic "github.com/GeorgeD19/libraries/json-logic-go"
	"github.com/buger/jsonparser"
	"github.com/h2non/filetype"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cast"
)

var (
	ErrInvalidOperation = errors.New("Invalid operation")
)

const (
	colCount = 2
	marginRw = 4.0
	marginH  = 15.0
	lineHtLg = 6.0
	lineHt   = 4.0
)

type ItemData struct {
	Key   string
	Field string
	Value string
	Media []alpaca.ImageFile
}

type Configuration struct {
	currentX int
	currentY int
}

// RunOperations will iterate through the array of operations and execute each
func RunOperations(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf, err error) {

	parseErr := err

	jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			pdf, err = ParseObject(pdf, string(value), data)
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
func RunOperation(pdf *gofpdf.Fpdf, name string, logic string, data []ItemData) (opdf *gofpdf.Fpdf, err error) {
	// fmt.Println(name)
	switch name {
	case "new":
		pdf = New(pdf, logic, data)
		break
	case "addpage":
		pdf = AddPage(pdf, logic, data)
		break
	case "setfont":
		pdf = SetFont(pdf, logic, data)
		break
	case "setx":
		pdf = SetX(pdf, logic, data)
		break
	case "sety":
		pdf = SetY(pdf, logic, data)
		break
	case "setxy":
		pdf = SetXY(pdf, logic, data)
		break
	case "cell":
		pdf = Cell(pdf, logic, data)
		break
	case "cellformat":
		pdf = CellFormat(pdf, logic, data)
		break
	case "setautopagebreak":
		pdf = SetAutoPageBreak(pdf, logic, data)
		break
	case "aliasnbpages":
		pdf = AliasNbPages(pdf, logic, data)
		break
	case "setheaderfunc":
		pdf, _ = SetHeaderFunc(pdf, logic, data)
		break
	case "setfooterfunc":
		pdf, _ = SetFooterFunc(pdf, logic, data)
		break
	case "settopmargin":
		pdf = SetTopMargin(pdf, logic, data)
		break
	case "setleftmargin":
		pdf = SetLeftMargin(pdf, logic, data)
		break
	case "settextcolor":
		pdf = SetTextColor(pdf, logic, data)
		break
	case "setfillcolor":
		pdf = SetFillColor(pdf, logic, data)
		break
	case "bodyfunc":
		pdf, _ = BodyFunc(pdf, logic, data)
		break
	case "ln":
		pdf = Ln(pdf, logic, data)
		break
	case "image":
		pdf = Image(pdf, logic, data)
		break

	case "var":
		pdf = Var(pdf, logic, data)
		break
	case "-":
		pdf = Minus(pdf, logic, data)
		break
	default:
		return pdf, ErrInvalidOperation
	}
	return pdf, nil
}

func Minus(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	// fmt.Println("Minus")
	// fmt.Println(logic)
	// TODO Pass in JSON values
	values := jsonlogic.GetValues(logic, `{"test":"4"}`)
	result := jsonlogic.MinusOperation(cast.ToFloat64(values[0]), cast.ToFloat64(values[1]))
	fmt.Println(values)
	fmt.Println(result)
	// pdf, _ = RunOperations(pdf, logic, data)

	return pdf
}

func Var(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	fmt.Println("Var")
	fmt.Println(logic) // jsonlogic.Apply()

	// docWidth, _ := pdf.GetPageSize()
	// docWidth = docWidth - (marginH * 2)
	// tr := pdf.UnicodeTranslatorFromDescriptor("")

	// fmt.Println(data)
	// // Try extract the variable from the data
	for x := 0; x < len(data); x++ {
		if data[x].Key == logic {
			fmt.Println(data[x].Value)
			// 		pdf.MultiCell(docWidth/2, lineHt, tr(data[x].Field), "", "L", true)
		}
	}

	return pdf
}

func Image(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	src := GetString("src", logic, "")
	name := GetString("name", logic, "")
	x := GetFloat("x", logic, 0.0)
	y := GetFloat("y", logic, 0.0)
	width := GetFloat("width", logic, 0.0)
	height := GetFloat("height", logic, 0.0)
	flow := GetBool("flow", logic, false)
	link := GetInt("link", logic, -1)
	linkStr := GetString("linkstr", logic, "")

	image, _ := GetImage(src)
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

func Ln(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	height := GetFloat("height", logic, -1.0)
	pdf.Ln(height)
	return pdf
}

func BodyFunc(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf, err error) {
	// fmt.Println("bodyfunc")
	// fmt.Println(data)
	for x := 0; x < len(data); x++ {
		item := make([]ItemData, 0)
		item = append(item, data[x])
		pdf, err = RunOperations(pdf, logic, item)
	}
	return pdf, err
}

func SetFillColor(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	r := GetInt("r", logic, 0)
	g := GetInt("r", logic, 0)
	b := GetInt("r", logic, 0)
	pdf.SetFillColor(r, g, b)
	return pdf
}

func SetTextColor(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	r := GetInt("r", logic, 0)
	g := GetInt("r", logic, 0)
	b := GetInt("r", logic, 0)
	pdf.SetTextColor(r, g, b)
	return pdf
}

func New(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	orientation := GetString("orientation", logic, "P")
	unit := GetString("unit", logic, "mm")
	size := GetString("size", logic, "A4")
	directory := GetString("dir", logic, "")
	return gofpdf.New(orientation, unit, size, directory)
}

func AddPage(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	pdf.AddPage()
	return pdf
}

func SetFont(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	family := GetString("family", logic, "Arial")
	style := GetString("style", logic, "")
	size := GetFloat("size", logic, 8.0)
	pdf.SetFont(family, style, size)
	return pdf
}

func SetAutoPageBreak(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	auto := GetBool("auto", logic, true)
	margin := GetFloat("margin", logic, 15.0)
	pdf.SetAutoPageBreak(auto, margin)
	return pdf
}

func AliasNbPages(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	aliasStr := GetString("alias", logic, "")
	pdf.AliasNbPages(aliasStr)
	return pdf
}

func Cell(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	width := GetFloat("width", logic, 0.0)
	height := GetFloat("height", logic, 0.0)
	text := GetString("text", logic, "")
	pdf.Cell(width, height, text)
	return pdf
}

func CellFormat(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	width := GetFloat("width", logic, 0.0)
	height := GetFloat("height", logic, 0.0)
	text := GetString("text", logic, "")
	text = strings.Replace(text, "{nn}", cast.ToString(pdf.PageNo()), -1)
	border := GetString("border", logic, "")
	line := GetInt("line", logic, 0)
	align := GetString("align", logic, "L")
	fill := GetBool("fill", logic, false)
	link := GetInt("link", logic, 0)
	linkStr := GetString("linkstr", logic, "")
	pdf.CellFormat(width, height, text, border, line, align, fill, link, linkStr)
	return pdf
}

func SetY(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	y := GetFloat("y", logic, 0.0)
	pdf.SetY(y)
	return pdf
}

func SetX(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	x := GetFloat("x", logic, 0.0)
	pdf.SetX(x)
	return pdf
}

func SetXY(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	x := GetFloat("x", logic, 0.0)
	y := GetFloat("y", logic, 0.0)
	pdf.SetXY(x, y)
	return pdf
}

func SetTopMargin(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	margin := GetFloat("margin", logic, 0.0)
	pdf.SetTopMargin(margin)
	return pdf
}

func SetLeftMargin(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf) {
	margin := GetFloat("margin", logic, 0.0)
	pdf.SetLeftMargin(margin)
	return pdf
}

func SetHeaderFunc(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf, err error) {
	pdf.SetHeaderFunc(func() {
		pdf, err = RunOperations(pdf, logic, data)
	})
	return pdf, err
}

func SetFooterFunc(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf, err error) {
	pdf.SetFooterFunc(func() {
		pdf, err = RunOperations(pdf, logic, data)
	})
	return pdf, err
}

// ParseObject entry point
func ParseObject(pdf *gofpdf.Fpdf, logic string, data []ItemData) (opdf *gofpdf.Fpdf, err error) {

	err = jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		pdf, err = RunOperation(pdf, string(key), string(value), data)
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

func GetBool(name string, logic string, fallback bool) (value bool) {
	result := fallback
	attribute, _, _, err := GetAttribute(name, logic, false)
	if err == nil {
		result = cast.ToBool(attribute)
	}
	return result
}

func GetFloat(name string, logic string, fallback float64) (value float64) {
	result := fallback
	attribute, _, _, err := GetAttribute(name, logic, false)
	if err == nil {
		result = cast.ToFloat64(cast.ToString(attribute))
	}
	return result
}

func GetInt(name string, logic string, fallback int) (value int) {
	result := fallback
	attribute, _, _, err := GetAttribute(name, logic, false)
	if err == nil {
		result = cast.ToInt(cast.ToString(attribute))
	}
	return result
}

func GetString(name string, logic string, fallback string) (value string) {
	result := fallback
	attribute, _, _, err := GetAttribute(name, logic, true)
	if err == nil {
		result = cast.ToString(attribute)
	}
	return result
}

func GetAttribute(name string, logic string, debug bool) (value []byte, dataType jsonparser.ValueType, offset int, err error) {

	value, dataType, offset, err = jsonparser.Get([]byte(logic), name)
	// TODO Determine if the getvalues has actually given something or if we should take the straight out value? or maybe that's already done as getvalues detects the type
	// attr := jsonlogic.GetValues(string(value), `{"test":"4"}`)
	// value = []byte(attr)
	// fmt.Println(cast.ToString(value))
	// fmt.Println(cast.ToInt(cast.ToString(attr)))

	// Here we can get the attribute from the object and if it matches any logic operations then we run them
	return value, dataType, offset, err
}

// Apply is the entry function to parse logic and optional data
func Apply(logic string, data []ItemData) (opdf *gofpdf.Fpdf, errs error) {
	pdf := new(gofpdf.Fpdf)
	result, err := RunOperations(pdf, logic, data)

	return result, err
}

type ImageFile struct {
	Data   string
	Width  int
	Height int
	Type   string
	Mime   string
}

// GetImage returns a File type containing the hex version of the image with associated meta information
func GetImage(FileName string) (f ImageFile, Err error) {

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
		FoundFile.Width, FoundFile.Height = GetJpgDimensions(File)
	case "gif":
		FoundFile.Width, FoundFile.Height = GetGifDimensions(File)
	case "png":
		FoundFile.Width, FoundFile.Height = GetPngDimensions(File)
	case "bmp":
		FoundFile.Width, FoundFile.Height = GetBmpDimensions(File)
	default:
		return FoundFile, err
	}

	return FoundFile, nil
}

func GetJpgDimensions(file *os.File) (width int, height int) {
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

func GetGifDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 4)
	file.ReadAt(bytes, 6)
	width = int(bytes[0]) + int(bytes[1])*256
	height = int(bytes[2]) + int(bytes[3])*256
	return
}

func GetBmpDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 18)
	width = int(bytes[3])<<24 | int(bytes[2])<<16 | int(bytes[1])<<8 | int(bytes[0])
	height = int(bytes[7])<<24 | int(bytes[6])<<16 | int(bytes[5])<<8 | int(bytes[4])
	return
}

func GetPngDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 16)
	width = int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	height = int(bytes[4])<<24 | int(bytes[5])<<16 | int(bytes[6])<<8 | int(bytes[7])
	return
}
