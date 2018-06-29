package jsongofpdf

import (
	jsonlogic "github.com/GeorgeD19/json-logic-go"
	"github.com/spf13/cast"
)

func (p *JSONGOFPDF) Calculation(logic, value string) string {
	result := 0.0
	calcType := p.GetString("type", logic, "")
	formula := p.GetString("formula", logic, "")
	table := p.Tables[p.TableIndex]
	switch calcType {
	case "count":
		for _, data := range table.Data {
			logicresult, _ := jsonlogic.Apply(formula, data)
			if cast.ToFloat64(logicresult) > 0 {
				result += 1.0
			}
		}
		break
	case "sum":
		for _, data := range table.Data {
			logicresult, _ := jsonlogic.Apply(formula, data)
			result += cast.ToFloat64(logicresult)
		}
		break
	case "minimum":
		for i, data := range table.Data {
			logicresult, _ := jsonlogic.Apply(formula, data)
			if i > 0 {
				if cast.ToFloat64(logicresult) < result {
					result = cast.ToFloat64(logicresult)
				}
			} else {
				result += cast.ToFloat64(logicresult)
			}
		}
		break
	case "maximum":
		for i, data := range table.Data {
			logicresult, _ := jsonlogic.Apply(formula, data)
			if i > 0 {
				if cast.ToFloat64(logicresult) > result {
					result = cast.ToFloat64(logicresult)
				}
			} else {
				result += cast.ToFloat64(logicresult)
			}
		}
		break
	case "average":
		for _, data := range table.Data {
			logicresult, _ := jsonlogic.Apply(formula, data)
			result += cast.ToFloat64(logicresult)
		}
		result = result / cast.ToFloat64(len(table.Data))
		break
	default:
		if len(table.Data)-1 >= p.RowIndex {
			logicresult, _ := jsonlogic.Apply(formula, table.Data[p.RowIndex])
			result = cast.ToFloat64(logicresult)
		}
		break
	}

	if result <= 0 && value != "" {
		return value
	}

	return cast.ToString(result)
}
