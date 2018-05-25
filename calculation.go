package jsongofpdf

import (
	"github.com/buger/jsonparser"
	"github.com/spf13/cast"
)

func (p *JSONGOFPDF) Calculation(logic, value string) string {
	result := value
	calcType := p.GetString("type", logic, "")
	switch calcType {
	case "count":

		break
	case "sum":
		result = p.CalculationSum(logic, value)
		break
	case "minimum":

		break
	case "maximum":

		break
	case "average":

		break
	}
	return result
}

func (p *JSONGOFPDF) CalculationSum(logic, value string) string {
	result := value
	formula := p.GetString("formula", logic, "")
	if formula != "" {
		calcResult := 0.0
		jsonparser.ArrayEach([]byte(formula), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			switch dataType {
			case jsonparser.String:

				// foreach value in formula array we get all the cell values with said key and add their float values up
				for _, row := range p.Tables[p.TableIndex].Rows {
					for _, cell := range row.Cells {
						target := string(value)
						if target == cell.Key || target == cell.Path {
							calcResult += cast.ToFloat64(cell.Value)
						}
					}
				}

				break
			}
		})
		result = cast.ToString(calcResult)
	}
	return result
}
