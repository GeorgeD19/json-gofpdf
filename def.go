package jsongofpdf

import (
	"errors"

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
	// Row options
	RowIndex       int
	RowHeight      float64
	RowCells       float64
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
}

type RowOptions struct {
	Index          int
	CurrentY       float64
	NextY          float64
	PrevCellHeight float64
}

type ImageFile struct {
	Data   string
	Width  int
	Height int
	Type   string
	Mime   string
}
