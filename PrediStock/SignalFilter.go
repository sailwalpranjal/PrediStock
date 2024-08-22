package mop

import (
	"strconv"
	"strings"
)

type Filter struct {
	profile *Profile
}

func NewFilter(profile *Profile) *Filter {
	return &Filter{
		profile: profile,
	}
}
func stringToNumber(numberString string) float64 {
	newString := strings.TrimSpace(numberString)
	newString = strings.Replace(newString, "$", "", 1)
	newString = strings.Replace(newString, "%", "", 1)
	newString = strings.Replace(newString, "K", "E+3", 1)
	newString = strings.Replace(newString, "M", "E+6", 1)
	newString = strings.Replace(newString, "B", "E+9", 1)
	newString = strings.Replace(newString, "T", "E+12", 1)
	finalValue, _ := strconv.ParseFloat(newString, 64)
	return finalValue
}
func (filter *Filter) Apply(stocks []Stock) []Stock {
	var filteredStocks []Stock

	for _, stock := range stocks {
		var values = make(map[string]interface{})
		values["ticker"] = strings.TrimSpace(stock.Ticker)
		values["last"] = stringToNumber(stock.LastTrade)
		values["change"] = stringToNumber(stock.Change)
		values["changePercent"] = stringToNumber(stock.ChangePct)
		values["open"] = stringToNumber(stock.Open)
		values["low"] = stringToNumber(stock.Low)
		values["high"] = stringToNumber(stock.High)
		values["low52"] = stringToNumber(stock.Low52)
		values["high52"] = stringToNumber(stock.High52)
		values["dividend"] = stringToNumber(stock.Dividend)
		values["yield"] = stringToNumber(stock.Yield)
		values["mktCap"] = stringToNumber(stock.MarketCap)
		values["mktCapX"] = stringToNumber(stock.MarketCapX)
		values["volume"] = stringToNumber(stock.Volume)
		values["avgVolume"] = stringToNumber(stock.AvgVolume)
		values["pe"] = stringToNumber(stock.PeRatio)
		values["peX"] = stringToNumber(stock.PeRatioX)
		values["direction"] = stock.Direction

		result, err := filter.profile.filterExpression.Evaluate(values)

		if err != nil {
			filter.profile.Filter = ""
			return filteredStocks
		}

		truthy, ok := result.(bool)

		if !ok {
			filter.profile.Filter = ""
			return filteredStocks
		}

		if truthy {
			filteredStocks = append(filteredStocks, stock)
		}
	}

	return filteredStocks
}
