package jsongofpdf

import (
	"strconv"

	"github.com/leekchan/accounting"
	"github.com/vjeantet/jodaTime"
)

func (p *JSONGOFPDF) Format(logic, value string) string {
	result := value
	formatType := p.GetString("type", logic, logic)
	switch formatType {
	case "currency":
		result = p.FormatCurrency(logic, value)
		break
	case "date":
		result = p.FormatDate(logic, value)
		break
	}
	return result
}

func (p *JSONGOFPDF) FormatCurrency(logic, value string) string {
	result := value
	resultFloat, err := strconv.ParseFloat(result, 64)
	if err == nil {
		symbol := p.GetString("symbol", logic, "£")
		precision := p.GetInt("precision", logic, 2)
		ac := accounting.Accounting{Symbol: symbol, Precision: precision}
		formatResult := ac.FormatMoney(resultFloat)
		if formatResult != "" {
			result = formatResult
		}
	}
	return result
}

func (p *JSONGOFPDF) FormatDate(logic, value string) string {
	result := value
	parse := p.GetString("parse", logic, "yyyy-M-d")
	format := p.GetString("format", logic, "d/M/yyyy")
	parseDate, err := jodaTime.Parse(parse, value)
	if err == nil {
		result = jodaTime.Format(format, parseDate)
	}
	return result
}
