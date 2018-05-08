package jsongofpdf

import (
	"errors"

	alpaca "github.com/GeorgeD19/alpaca-go"
	"github.com/GeorgeD19/backup/repository/helper"
	"github.com/GeorgeD19/backup/repository/model"
)

var (
	ErrDefaultError     = errors.New("You must supply at least one argument.")
	ErrInvalidOperation = errors.New("Invalid operation")
)

type JSONGOFPDF struct {
	Logic      string
	Data       string
	Order      *model.Order
	Submission *model.Submission
	Form       *model.Form
	Parser     *alpaca.Alpaca
	User       helper.User
	ColCount   int
	MarginRw   float64
	MarginH    float64
	LineHtLg   float64
	LineHt     float64
	DocWidth   float64
	initY      float64
	initX      float64
	currentX   float64
	currentY   float64
	nextY      float64
	tr         func(string) string
}

type JSONGOFPDFOptions struct {
	Logic      string
	Data       string
	Order      *model.Order
	Submission *model.Submission
	Form       *model.Form
	Parser     *alpaca.Alpaca
	User       helper.User
}
