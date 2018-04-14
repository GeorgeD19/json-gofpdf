package jsongofpdf

import (
	"testing"
)

// func TestBlankPage(t *testing.T) {
// 	logic := `[
// 			{
// 				"New": {
// 					"orientation": "P",
// 					"unit": "mm",
// 					"size": "A4",
// 					"fontDir": ""
// 				}
// 			},
// 			{
// 				"SetFont": {
// 					"family": "Arial",
// 					"style": "",
// 					"size": 8.9
// 				}
// 			},
// 			{
// 				"SetAutoPageBreak": {
// 					"auto": true,
// 					"margin": 15.0
// 				}
// 			}
// 		]`

// 	data := `{}`

// 	// Should return true
// 	result, _, _ := Apply(logic, data)

// 	if result != "P" {
// 		t.Fatal("Logic should return P")
// 	}
// }

func TestGetString(t *testing.T) {
	logic := `
	"new": {
		"orientation": "P"
	}`

	// Should return true
	result := GetString("orientation", logic, "")

	if result != "P" {
		t.Fatal("Logic should return P")
	}
}

func TestGetInt(t *testing.T) {
	logic := `
	"setfont": {
		"size": 10
	}`

	// Should return true
	result := GetInt("size", logic, 8)

	if result != 10 {
		t.Fatal("Logic should return 10")
	}
}

func TestGetStringFallback(t *testing.T) {
	logic := `"new": {}`

	// Should return true
	result := GetString("orientation", logic, "P")

	if result != "P" {
		t.Fatal("Logic should return P fallback")
	}
}
