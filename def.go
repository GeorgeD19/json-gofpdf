package jsongofpdf

import (
	"errors"
	"time"

	alpaca "github.com/GeorgeD19/alpaca-go"
	"github.com/GeorgeD19/securigroupgo/repository/model"
)

var (
	ErrDefaultError     = errors.New("You must supply at least one argument.")
	ErrInvalidOperation = errors.New("Invalid operation")
)

type JSONGOFPDF struct {
	Logic        string
	Data         string
	Order        *model.Order
	Submission   *model.Submission
	Form         *model.Form
	Parser       *alpaca.Alpaca
	User         model.User
	ColCount     int
	MarginRw     float64
	MarginH      float64
	LineHtLg     float64
	LineHt       float64
	DocWidth     float64
	initY        float64
	initX        float64
	nextY        float64
	tr           func(string) string
	DataIndex    int
	NewPage      bool
	currentPage  int
	HeaderHeight float64

	Globals map[string]string
	// Table options
	Tables     []Table
	TableIndex int
	// Row options
	RowFuncIndex int
	RowIndex     int
	RowHeight    float64
	RowCells     float64
	// Cell options
	CellIndex      int
	CurrentX       float64
	CurrentY       float64
	CurrentRowX    float64
	CurrentRowY    float64
	NextY          float64
	PrevCellHeight float64
	MediaIndex     int
	DPI            int
}

type JSONGOFPDFOptions struct {
	Logic      string
	Data       string
	Order      *model.Order
	Submission *model.Submission
	Form       *model.Form
	Parser     *alpaca.Alpaca
	User       model.User
	Tables     []Table
	Globals    map[string]string
}

type RowOptions struct {
	Index          int
	CurrentY       float64
	NextY          float64
	PrevCellHeight float64
}

type File struct {
	Data    string
	Type    string
	Mime    string
	Name    string
	Created time.Time
}

type ImageFile struct {
	Data   string
	Type   string
	Mime   string
	Name   string
	Width  int
	Height int
}

type Table struct {
	Rows []Row
}

type Row struct {
	Cells []Cell
}

type Cell struct {
	Path   string
	Key    string
	Title  string
	Value  interface{}
	Images []ImageFile
}
