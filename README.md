# json-gofpdf

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/GeorgeD19/jsPDF-json/blob/master/license.txt)

json-gofpdf runs JSON schema through [gofpdf](https://github.com/jung-kurt/gofpdf) to generate PDFs in server-side.

## Features
- Maps JSON schema to gofpdf functions.
- Inject data into your schema layouts for dynamic information.

## API

```golang
	logic := `[
		{
			"new": {
				"orientation": "P",
				"unit": "mm",
				"size": "A4",
				"fontDir": ""
			}
		},
		{
			"addpage": {}
		},
		{
			"setfont": {}
		},
		{
			"cell": {
				"width": 40.0,
				"height": 10.0,
				"text": "Hello, World!"
			}
		}
	]`
	
	pdfparser2, err := jsongofpdf.New(jsongofpdf.JSONGOFPDFOptions{Logic: logic})
	if err != nil {
		fmt.Println(err)
	}
	pdf2 := pdfparser2.GetPDF()
	pdf.OutputFileAndClose("hello.pdf")
```

### Supported functions
- AddPage
- AliasNbPages
- Body
- Cell
- CellFormat
- 