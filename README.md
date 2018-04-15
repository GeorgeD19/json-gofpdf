# json-gofpdf
Parse json and run through gofpdf

```
pdfLogic := `[
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

	data := `{}`
  pdf, _, _ := jsonpdf.Apply(pdfLogic, data)
	pdf.OutputFileAndClose("hello.pdf")
  ```
