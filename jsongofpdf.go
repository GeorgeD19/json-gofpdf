package jsongofpdf

import (
	"errors"

	"github.com/buger/jsonparser"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cast"
)

var (
	ErrInvalidOperation = errors.New("Invalid Operation")
)

// RunOperations will iterate through the array of operations and execute each
func RunOperations(logic string, data string) (opdf *gofpdf.Fpdf, err error) {
	pdf := new(gofpdf.Fpdf)

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
func RunOperation(pdf *gofpdf.Fpdf, name string, logic string, data string) (opdf *gofpdf.Fpdf, err error) {
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
	case "cell":
		pdf = Cell(pdf, logic, data)
		break
	default:
		return pdf, ErrInvalidOperation
	}
	return pdf, nil
}

func New(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf) {
	orientation := GetString("orientation", logic, "P")
	unit := GetString("unit", logic, "mm")
	size := GetString("size", logic, "A4")
	directory := GetString("dir", logic, "")
	return gofpdf.New(orientation, unit, size, directory)
}

func AddPage(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf) {
	pdf.AddPage()
	return pdf
}

func SetFont(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf) {
	family := GetString("family", logic, "Arial")
	style := GetString("style", logic, "")
	size := GetFloat("size", logic, 8.0)
	pdf.SetFont(family, style, size)
	return pdf
}

func Cell(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf) {
	width := GetFloat("width", logic, 0.0)
	height := GetFloat("height", logic, 0.0)
	text := GetString("text", logic, "")
	pdf.Cell(width, height, text)
	return pdf
}

// ParseObject entry point
func ParseObject(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf, err error) {

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
	attribute, _, _, err := GetAttribute(name, logic)
	if err == nil {
		result = cast.ToBool(attribute)
	}
	return result
}

func GetFloat(name string, logic string, fallback float64) (value float64) {
	result := fallback
	attribute, _, _, err := GetAttribute(name, logic)
	if err == nil {
		result = cast.ToFloat64(cast.ToString(attribute))
	}
	return result
}

func GetInt(name string, logic string, fallback int) (value int) {
	result := fallback
	attribute, _, _, err := GetAttribute(name, logic)
	if err == nil {
		result = cast.ToInt(cast.ToString(attribute))
	}
	return result
}

func GetString(name string, logic string, fallback string) (value string) {
	result := fallback
	attribute, _, _, err := GetAttribute(name, logic)
	if err == nil {
		result = cast.ToString(attribute)
	}
	return result
}

func GetAttribute(name string, logic string) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	return jsonparser.Get([]byte(logic), name)
}

// Apply is the entry function to parse logic and optional data
func Apply(logic string, data string) (opdf *gofpdf.Fpdf, errs error) {

	// Ensure data is object
	if data == `` {
		data = `{}`
	}

	// Must be an object to kick off process
	result, err := RunOperations(logic, data)

	return result, err
}
