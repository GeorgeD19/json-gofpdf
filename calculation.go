package jsongofpdf

import (
	jsonlogic "github.com/GeorgeD19/json-logic-go"
	"github.com/spf13/cast"
)

func (p *JSONGOFPDF) Calculation(logic, value string) string {
	result := 0.0
	calcType := p.GetString("type", logic, "")
	formula := p.GetString("formula", logic, "")
	switch calcType {
	case "count":
		for _, data := range p.Tables[p.TableIndex].Data {
			logicresult, _ := jsonlogic.Apply(formula, data)
			if cast.ToFloat64(logicresult) > 0 {
				result += 1.0
			}
		}
		break
	case "sum":
		for _, data := range p.Tables[p.TableIndex].Data {
			logicresult, _ := jsonlogic.Apply(formula, data)
			result += cast.ToFloat64(logicresult)
		}
		break
	case "minimum":
		for i, data := range p.Tables[p.TableIndex].Data {
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
		for i, data := range p.Tables[p.TableIndex].Data {
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
		for _, data := range p.Tables[p.TableIndex].Data {
			logicresult, _ := jsonlogic.Apply(formula, data)
			result += cast.ToFloat64(logicresult)
		}
		result = result / cast.ToFloat64(len(p.Tables[p.TableIndex].Data))
		break
	default:
		data := p.Tables[p.TableIndex].Data[p.TableIndex]
		logicresult, _ := jsonlogic.Apply(formula, data)
		result = cast.ToFloat64(logicresult)
		break
	}

	if result <= 0 && value != "" {
		return value
	}

	return cast.ToString(result)
}
