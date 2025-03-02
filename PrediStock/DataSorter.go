package mop

/*
The code defines a sorting system for stock data based on multiple attributes (e.g., ticker, last trade, change, volume, market cap, etc.). 
It supports both ascending and descending sorting orders, and the sorting criteria are driven by the profile's configuration.

- `Sorter` struct: Manages sorting behavior by interacting with the `Profile`.
- `sortable` type: Defines the list of stocks to be sorted.
- Sorting types (e.g., `byTickerAsc`, `byVolumeDesc`): Each represents a sorting criterion (ascending or descending) for different stock attributes.
- Sorting functions (`Less`): Compare stock attributes for sorting in the desired order (either ascending or descending).
- `NewSorter`: Initializes the `Sorter` with the user's profile settings.
- `SortByCurrentColumn`: Applies the appropriate sorting strategy based on the current profile settings (ascending/descending).
- Helper functions:
  - `c`: Converts and normalizes string values (e.g., percentage or monetary values).
  - `m`: Converts string representations of values with multipliers (K, M, B, T) into numerical values.
*/

import (
	"sort"
	"strconv"
	"strings"
)

type Sorter struct {
	profile *Profile
}

type sortable []Stock

func (list sortable) Len() int      { return len(list) }
func (list sortable) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

type byTickerAsc struct{ sortable }
type byLastTradeAsc struct{ sortable }
type byChangeAsc struct{ sortable }
type byChangePctAsc struct{ sortable }
type byOpenAsc struct{ sortable }
type byLowAsc struct{ sortable }
type byHighAsc struct{ sortable }
type byLow52Asc struct{ sortable }
type byHigh52Asc struct{ sortable }
type byVolumeAsc struct{ sortable }
type byAvgVolumeAsc struct{ sortable }
type byPeRatioAsc struct{ sortable }
type byDividendAsc struct{ sortable }
type byYieldAsc struct{ sortable }
type byMarketCapAsc struct{ sortable }
type byPreOpenAsc struct{ sortable }
type byAfterHoursAsc struct{ sortable }

type byTickerDesc struct{ sortable }
type byLastTradeDesc struct{ sortable }
type byChangeDesc struct{ sortable }
type byChangePctDesc struct{ sortable }
type byOpenDesc struct{ sortable }
type byLowDesc struct{ sortable }
type byHighDesc struct{ sortable }
type byLow52Desc struct{ sortable }
type byHigh52Desc struct{ sortable }
type byVolumeDesc struct{ sortable }
type byAvgVolumeDesc struct{ sortable }
type byPeRatioDesc struct{ sortable }
type byDividendDesc struct{ sortable }
type byYieldDesc struct{ sortable }
type byMarketCapDesc struct{ sortable }
type byPreOpenDesc struct{ sortable }
type byAfterHoursDesc struct{ sortable }

func (list byTickerAsc) Less(i, j int) bool {
	return list.sortable[i].Ticker < list.sortable[j].Ticker
}
func (list byLastTradeAsc) Less(i, j int) bool {
	return list.sortable[i].LastTrade < list.sortable[j].LastTrade
}
func (list byChangeAsc) Less(i, j int) bool {
	return c(list.sortable[i].Change) < c(list.sortable[j].Change)
}
func (list byChangePctAsc) Less(i, j int) bool {
	return c(list.sortable[i].ChangePct) < c(list.sortable[j].ChangePct)
}
func (list byOpenAsc) Less(i, j int) bool {
	return list.sortable[i].Open < list.sortable[j].Open
}
func (list byLowAsc) Less(i, j int) bool {
	return list.sortable[i].Low < list.sortable[j].Low
}
func (list byHighAsc) Less(i, j int) bool {
	return list.sortable[i].High < list.sortable[j].High
}
func (list byLow52Asc) Less(i, j int) bool {
	return list.sortable[i].Low52 < list.sortable[j].Low52
}
func (list byHigh52Asc) Less(i, j int) bool {
	return list.sortable[i].High52 < list.sortable[j].High52
}
func (list byVolumeAsc) Less(i, j int) bool {
	return m(list.sortable[i].Volume) < m(list.sortable[j].Volume)
}
func (list byAvgVolumeAsc) Less(i, j int) bool {
	return m(list.sortable[i].AvgVolume) < m(list.sortable[j].AvgVolume)
}
func (list byPeRatioAsc) Less(i, j int) bool {
	return list.sortable[i].PeRatio < list.sortable[j].PeRatio
}
func (list byDividendAsc) Less(i, j int) bool {
	return list.sortable[i].Dividend < list.sortable[j].Dividend
}
func (list byYieldAsc) Less(i, j int) bool {
	return list.sortable[i].Yield < list.sortable[j].Yield
}
func (list byMarketCapAsc) Less(i, j int) bool {
	return m(list.sortable[i].MarketCap) < m(list.sortable[j].MarketCap)
}
func (list byPreOpenAsc) Less(i, j int) bool {
	return c(list.sortable[i].PreOpen) < c(list.sortable[j].PreOpen)
}
func (list byAfterHoursAsc) Less(i, j int) bool {
	return c(list.sortable[i].AfterHours) < c(list.sortable[j].AfterHours)
}

func (list byTickerDesc) Less(i, j int) bool {
	return list.sortable[j].Ticker < list.sortable[i].Ticker
}
func (list byLastTradeDesc) Less(i, j int) bool {
	return list.sortable[j].LastTrade < list.sortable[i].LastTrade
}
func (list byChangeDesc) Less(i, j int) bool {
	return c(list.sortable[j].Change) < c(list.sortable[i].Change)
}
func (list byChangePctDesc) Less(i, j int) bool {
	return c(list.sortable[j].ChangePct) < c(list.sortable[i].ChangePct)
}
func (list byOpenDesc) Less(i, j int) bool {
	return list.sortable[j].Open < list.sortable[i].Open
}
func (list byLowDesc) Less(i, j int) bool {
	return list.sortable[j].Low < list.sortable[i].Low
}
func (list byHighDesc) Less(i, j int) bool {
	return list.sortable[j].High < list.sortable[i].High
}
func (list byLow52Desc) Less(i, j int) bool {
	return list.sortable[j].Low52 < list.sortable[i].Low52
}
func (list byHigh52Desc) Less(i, j int) bool {
	return list.sortable[j].High52 < list.sortable[i].High52
}
func (list byVolumeDesc) Less(i, j int) bool {
	return m(list.sortable[j].Volume) < m(list.sortable[i].Volume)
}
func (list byAvgVolumeDesc) Less(i, j int) bool {
	return m(list.sortable[j].AvgVolume) < m(list.sortable[i].AvgVolume)
}
func (list byPeRatioDesc) Less(i, j int) bool {
	return list.sortable[j].PeRatio < list.sortable[i].PeRatio
}
func (list byDividendDesc) Less(i, j int) bool {
	return list.sortable[j].Dividend < list.sortable[i].Dividend
}
func (list byYieldDesc) Less(i, j int) bool {
	return list.sortable[j].Yield < list.sortable[i].Yield
}
func (list byMarketCapDesc) Less(i, j int) bool {
	return m(list.sortable[j].MarketCap) < m(list.sortable[i].MarketCap)
}
func (list byPreOpenDesc) Less(i, j int) bool {
	return c(list.sortable[j].PreOpen) < c(list.sortable[i].PreOpen)
}
func (list byAfterHoursDesc) Less(i, j int) bool {
	return c(list.sortable[j].AfterHours) < c(list.sortable[i].AfterHours)
}
func NewSorter(profile *Profile) *Sorter {
	return &Sorter{
		profile: profile,
	}
}
func (sorter *Sorter) SortByCurrentColumn(stocks []Stock) *Sorter {
	var interfaces []sort.Interface

	if sorter.profile.Ascending {
		interfaces = []sort.Interface{
			byTickerAsc{stocks},
			byLastTradeAsc{stocks},
			byChangeAsc{stocks},
			byChangePctAsc{stocks},
			byOpenAsc{stocks},
			byLowAsc{stocks},
			byHighAsc{stocks},
			byLow52Asc{stocks},
			byHigh52Asc{stocks},
			byVolumeAsc{stocks},
			byAvgVolumeAsc{stocks},
			byPeRatioAsc{stocks},
			byDividendAsc{stocks},
			byYieldAsc{stocks},
			byMarketCapAsc{stocks},
			byPreOpenAsc{stocks},
			byAfterHoursAsc{stocks},
		}
	} else {
		interfaces = []sort.Interface{
			byTickerDesc{stocks},
			byLastTradeDesc{stocks},
			byChangeDesc{stocks},
			byChangePctDesc{stocks},
			byOpenDesc{stocks},
			byLowDesc{stocks},
			byHighDesc{stocks},
			byLow52Desc{stocks},
			byHigh52Desc{stocks},
			byVolumeDesc{stocks},
			byAvgVolumeDesc{stocks},
			byPeRatioDesc{stocks},
			byDividendDesc{stocks},
			byYieldDesc{stocks},
			byMarketCapDesc{stocks},
			byPreOpenDesc{stocks},
			byAfterHoursDesc{stocks},
		}
	}

	sort.Sort(interfaces[sorter.profile.SortColumn])

	return sorter
}
func c(str string) float32 {
	c := "$"
	for _, v := range currencies {
		if strings.Contains(str, v) {
			c = v
		}
	}
	trimmed := strings.Replace(strings.Trim(str, ` %`), c, ``, 1)
	value, _ := strconv.ParseFloat(trimmed, 32)
	return float32(value)
}
func m(str string) float32 {
	if len(str) == 0 {
		return 0
	}
	multiplier := 1.0
	switch str[len(str)-1:] {
	case `T`:
		multiplier = 1000000000000.0
	case `B`:
		multiplier = 1000000000.0
	case `M`:
		multiplier = 1000000.0
	case `K`:
		multiplier = 1000.0
	}
	trimmed := strings.Trim(str, ` $TBMK`)
	value, _ := strconv.ParseFloat(trimmed, 32)
	return float32(value * multiplier)
}
