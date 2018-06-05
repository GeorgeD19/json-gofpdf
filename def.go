package jsongofpdf

import (
	"errors"
	"time"
)

var (
	ErrDefaultError     = errors.New("You must supply at least one argument.")
	ErrInvalidOperation = errors.New("Invalid operation")
)

type JSONGOFPDFOptions struct {
	Logic   string
	Data    string
	Tables  []Table
	Globals map[string]interface{}
}

type JSONGOFPDF struct {
	Globals  map[string]interface{}
	tr       func(string) string
	Logic    string
	DocWidth float64
	initY    float64

	// Table options
	Tables     []Table
	TableIndex int
	// Row options
	RowFuncIndex int
	RowIndex     int
	RowHeight    float64
	RowCells     float64
	// Cell options
	CellPreIndex int
	CellIndex    int
	CurrentX     float64
	CurrentY     float64
	CurrentRowX  float64
	CurrentRowY  float64
	NextY        float64
	// Media options
	MediaIndex int
	DPI        int
}

type Table struct {
	Rows []Row
	Data []string
}

type Row struct {
	Cells []Cell
}

type Cell struct {
	Path     string
	Key      string
	Format   string
	Title    string
	Value    interface{}
	Images   []ImageFile
	Disabled bool
	Type     string
}

type ImageFile struct {
	Data    string
	Type    string
	Mime    string
	Name    string
	Created time.Time
	Width   int
	Height  int
}
