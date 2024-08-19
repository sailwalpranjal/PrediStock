package mop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const quotesURL = `https://query1.finance.yahoo.com/v7/finance/quote?crumb=%s&symbols=%s`

const quotesURLQueryParts = `&range=1d&interval=5m&indicators=close&includeTimestamps=false&includePrePost=false&corsDomain=finance.yahoo.com&.tsrc=finance`

const noDataIndicator = `N/A`

type Stock struct {
	Ticker     string `json:"symbol"`                      // Stock ticker.
	LastTrade  string `json:"regularMarketPrice"`          // l1: last trade.
	Change     string `json:"regularMarketChange"`         // c6: change real time.
	ChangePct  string `json:"regularMarketChangePercent"`  // k2: percent change real time.
	Open       string `json:"regularMarketOpen"`           // o: market open price.
	Low        string `json:"regularMarketDayLow"`         // g: day's low.
	High       string `json:"regularMarketDayHigh"`        // h: day's high.
	Low52      string `json:"fiftyTwoWeekLow"`             // j: 52-weeks low.
	High52     string `json:"fiftyTwoWeekHigh"`            // k: 52-weeks high.
	Volume     string `json:"regularMarketVolume"`         // v: volume.
	AvgVolume  string `json:"averageDailyVolume10Day"`     // a2: average volume.
	PeRatio    string `json:"trailingPE"`                  // r2: P/E ration real time.
	PeRatioX   string `json:"trailingPE"`                  // r: P/E ration (fallback when real time is N/A).
	Dividend   string `json:"trailingAnnualDividendRate"`  // d: dividend.
	Yield      string `json:"trailingAnnualDividendYield"` // y: dividend yield.
	MarketCap  string `json:"marketCap"`                   // j3: market cap real time.
	MarketCapX string `json:"marketCap"`                   // j1: market cap (fallback when real time is N/A).
	Currency   string `json:"currency"`                    // String code for currency of stock.
	Direction  int    // -1 when change is < $0, 0 when change is = $0, 1 when change is > $0.
	PreOpen    string `json:"preMarketChangePercent,omitempty"`
	AfterHours string `json:"postMarketChangePercent,omitempty"`
}

type Quotes struct {
	market  *Market  // Pointer to Market.
	profile *Profile // Pointer to Profile.
	stocks  []Stock  // Array of stock quote data.
	errors  string   // Error string if any.
}

func NewQuotes(market *Market, profile *Profile) *Quotes {
	return &Quotes{
		market:  market,
		profile: profile,
		errors:  ``,
	}
}

func (quotes *Quotes) Fetch() (self *Quotes) {
	self = quotes
	if quotes.isReady() {
		defer func() {
			if err := recover(); err != nil {
				quotes.errors = fmt.Sprintf("\n\n\n\nError fetching stock quotes...\n%s", err)
			} else {
				quotes.errors = ""
			}
		}()

		url := fmt.Sprintf(quotesURL, quotes.market.crumb, strings.Join(quotes.profile.Tickers, `,`))

		client := http.Client{}
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		request.Header = http.Header{
			"Accept":          {"*/*"},
			"Accept-Language": {"en-US,en;q=0.5"},
			"Connection":      {"keep-alive"},
			"Content-Type":    {"application/json"},
			"Cookie":          {quotes.market.cookies},
			"Host":            {"query1.finance.yahoo.com"},
			"Origin":          {"https://finance.yahoo.com"},
			"Referer":         {"https://finance.yahoo.com"},
			"Sec-Fetch-Dest":  {"empty"},
			"Sec-Fetch-Mode":  {"cors"},
			"Sec-Fetch-Site":  {"same-site"},
			"TE":              {"trailers"},
			"User-Agent":      {userAgent},
		}

		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		quotes.parse2(body)
	}

	return quotes
}
func (quotes *Quotes) Ok() (bool, string) {
	return quotes.errors == ``, quotes.errors
}
func (quotes *Quotes) AddTickers(tickers []string) (added int, err error) {
	if added, err = quotes.profile.AddTickers(tickers); err == nil && added > 0 {
		quotes.stocks = nil
	}
	return
}
func (quotes *Quotes) RemoveTickers(tickers []string) (removed int, err error) {
	if removed, err = quotes.profile.RemoveTickers(tickers); err == nil && removed > 0 {
		quotes.stocks = nil
	}
	return
}
func (quotes *Quotes) isReady() bool {
	return (quotes.stocks == nil || !quotes.market.IsClosed) && len(quotes.profile.Tickers) > 0
}
func (quotes *Quotes) parse2(body []byte) (*Quotes, error) {
	d := map[string]map[string][]map[string]interface{}{}
	err := json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	results := d["quoteResponse"]["result"]

	quotes.stocks = make([]Stock, len(results))
	for i, raw := range results {
		result := map[string]string{}
		for k, v := range raw {
			switch v.(type) {
			case string:
				result[k] = v.(string)
			case float64:
				result[k] = float2Str(v.(float64))
			default:
				result[k] = fmt.Sprintf("%v", v)
			}
		}
		quotes.stocks[i].Ticker = result["symbol"]
		quotes.stocks[i].LastTrade = result["regularMarketPrice"]
		quotes.stocks[i].Change = result["regularMarketChange"]
		quotes.stocks[i].ChangePct = result["regularMarketChangePercent"]
		quotes.stocks[i].Open = result["regularMarketOpen"]
		quotes.stocks[i].Low = result["regularMarketDayLow"]
		quotes.stocks[i].High = result["regularMarketDayHigh"]
		quotes.stocks[i].Low52 = result["fiftyTwoWeekLow"]
		quotes.stocks[i].High52 = result["fiftyTwoWeekHigh"]
		quotes.stocks[i].Volume = result["regularMarketVolume"]
		quotes.stocks[i].AvgVolume = result["averageDailyVolume10Day"]
		quotes.stocks[i].PeRatio = result["trailingPE"]
		quotes.stocks[i].PeRatioX = result["trailingPE"]
		quotes.stocks[i].Dividend = result["trailingAnnualDividendRate"]
		quotes.stocks[i].Yield = result["trailingAnnualDividendYield"]
		quotes.stocks[i].MarketCap = result["marketCap"]
		quotes.stocks[i].MarketCapX = result["marketCap"]
		quotes.stocks[i].Currency = result["currency"]
		quotes.stocks[i].PreOpen = result["preMarketChangePercent"]
		quotes.stocks[i].AfterHours = result["postMarketChangePercent"]
		adv, err := strconv.ParseFloat(quotes.stocks[i].Change, 64)
		quotes.stocks[i].Direction = 0
		if err == nil {
			if adv < 0.0 {
				quotes.stocks[i].Direction = -1
			} else if adv > 0.0 {
				quotes.stocks[i].Direction = 1
			}
		}
	}
	return quotes, nil
}
func (quotes *Quotes) parse(body []byte) *Quotes {
	lines := bytes.Split(body, []byte{'\n'})
	quotes.stocks = make([]Stock, len(lines))
	fieldsCount := reflect.ValueOf(quotes.stocks[0]).NumField() - 1
	for i, line := range lines {
		columns := bytes.Split(bytes.TrimSpace(line), []byte{','})
		for j := 0; j < fieldsCount; j++ {
			reflect.ValueOf(&quotes.stocks[i]).Elem().Field(j).SetString(string(columns[j]))
		}
		if quotes.stocks[i].PeRatio == `N/A` && quotes.stocks[i].PeRatioX != `N/A` {
			quotes.stocks[i].PeRatio = quotes.stocks[i].PeRatioX
		}
		if quotes.stocks[i].MarketCap == `N/A` && quotes.stocks[i].MarketCapX != `N/A` {
			quotes.stocks[i].MarketCap = quotes.stocks[i].MarketCapX
		}
		adv, err := strconv.ParseFloat(quotes.stocks[i].Change, 64)
		quotes.stocks[i].Direction = 0
		if err == nil {
			if adv < 0 {
				quotes.stocks[i].Direction = -1
			} else if adv > 0 {
				quotes.stocks[i].Direction = 1
			}
		}
	}

	return quotes
}

// -----------------------------------------------------------------------------
func sanitize(body []byte) []byte {
	return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}

func float2Str(v float64) string {
	unit := ""
	switch {
	case v > 1.0e12:
		v = v / 1.0e12
		unit = "T"
	case v > 1.0e9:
		v = v / 1.0e9
		unit = "B"
	case v > 1.0e6:
		v = v / 1.0e6
		unit = "M"
	case v > 1.0e5:
		v = v / 1.0e3
		unit = "K"
	default:
		unit = ""
	}
	return fmt.Sprintf("%0.3f%s", v, unit)
}
