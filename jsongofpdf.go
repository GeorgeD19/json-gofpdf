package jsongofpdf

import (
	"errors"

	"github.com/spf13/cast"

	"github.com/buger/jsonparser"
	"github.com/jung-kurt/gofpdf"
)

// Errors
var (
	ErrInvalidOperation = errors.New("Invalid Operation %s")
)

// Operations contains all possible operations that can be performed
var Operations = make(map[string]Operation)

// Operation interface that allows operations to be registered in a list
type Operation interface {
	run(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf, ovar interface{})
}

// AddOperation adds possible operation to Operations library
func AddOperation(name string, callable Operation) {
	Operations[name] = callable
}

// RemoveOperation removes possible operation from Operations library if it exists
func RemoveOperation(name string) {
	_, ok := Operations[name]
	if ok {
		delete(Operations, name)
	}
}

// RunOperation is to ensure that any operation ran doesn't crash the system if it doesn't exist
func RunOperation(pdf *gofpdf.Fpdf, name string, logic string, data string) (opdf *gofpdf.Fpdf, ores interface{}, err error) {
	_, ok := Operations[name]
	if ok {
		pdf, vars := Operations[name].run(pdf, logic, data)
		return pdf, vars, nil
	}
	return nil, nil, ErrInvalidOperation
}

func init() {
	AddOperation("new", new(New))
	AddOperation("addpage", new(AddPage))
	AddOperation("cell", new(Cell))
	AddOperation("setfont", new(SetFont))
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

// New type is entry point for parser
type New struct{}

func (o New) run(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf, ovar interface{}) {
	orientation := GetString("orientation", logic, "P")
	unit := GetString("unit", logic, "mm")
	size := GetString("size", logic, "A4")
	directory := GetString("dir", logic, "")
	return gofpdf.New(orientation, unit, size, directory), nil
}

// AddPage type is entry point for parser
type AddPage struct{}

func (o AddPage) run(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf, ovar interface{}) {
	pdf.AddPage()
	return pdf, nil
}

// SetFont type is entry point for parser
type SetFont struct{}

func (o SetFont) run(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf, ovar interface{}) {
	family := GetString("family", logic, "Arial")
	style := GetString("style", logic, "")
	size := GetFloat("size", logic, 8.0)
	pdf.SetFont(family, style, size)
	return pdf, nil
}

// Cell type is entry point for parser
type Cell struct{}

func (o Cell) run(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf, ovar interface{}) {
	width := GetFloat("width", logic, 0.0)
	height := GetFloat("height", logic, 0.0)
	text := GetString("text", logic, "")
	pdf.Cell(width, height, text)
	return pdf, nil
}

// RunOperations will iterate through the array of operations and execute each
func RunOperations(logic string, data string) (opdf *gofpdf.Fpdf) {
	pdf := new(gofpdf.Fpdf)
	jsonparser.ArrayEach([]byte(logic), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Object:
			pdf, _, _ = ParseObject(pdf, string(value), data)
			break
		}
	})

	return pdf
}

// ParseObject entry point
func ParseObject(pdf *gofpdf.Fpdf, logic string, data string) (opdf *gofpdf.Fpdf, vars interface{}, err error) {

	err = jsonparser.ObjectEach([]byte(logic), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if operation, ok := Operations[string(key)]; ok {
			pdf, _ = operation.run(pdf, string(value), data)
		} else {
			return ErrInvalidOperation
		}
		return nil
	})

	if err != nil {
		return nil, false, err
	}

	return pdf, vars, nil
}

// Apply is the entry function to parse logic and optional data
func Apply(logic string, data string) (opdf *gofpdf.Fpdf, vars interface{}, errs error) {

	// Ensure data is object
	if data == `` {
		data = `{}`
	}

	// Must be an object to kick off process
	result := RunOperations(logic, data)

	return result, nil, nil
}
