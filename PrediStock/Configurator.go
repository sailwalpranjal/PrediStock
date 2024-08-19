package mop

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/Knetic/govaluate"
)

const defaultGainColor = "green"
const defaultLossColor = "red"
const defaultTagColor = "yellow"
const defaultHeaderColor = "lightgray"
const defaultTimeColor = "lightgray"
const defaultColor = "lightgray"

type Profile struct {
	Tickers       []string // List of stock tickers to display.
	MarketRefresh int      // Time interval to refresh market data.
	QuotesRefresh int      // Time interval to refresh stock quotes.
	SortColumn    int      // Column number by which we sort stock quotes.
	Ascending     bool     // True when sort order is ascending.
	Grouped       bool     // True when stocks are grouped by advancing/declining.
	Filter        string   // Filter in human form
	UpDownJump    int      // Number of lines to go up/down when scrolling.
	Colors        struct { // User defined colors
		Gain    string
		Loss    string
		Tag     string
		Header  string
		Time    string
		Default string
	}
	ShowTimestamp    bool                          
	filterExpression *govaluate.EvaluableExpression 
	selectedColumn   int                           
	filename         string                        
}

func IsSupportedColor(colorName string) bool {
	switch colorName {
	case
		"black",
		"red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		"white",
		"darkgray",
		"lightred",
		"lightgreen",
		"lightyellow",
		"lightblue",
		"lightmagenta",
		"lightcyan",
		"lightgray":
		return true
	}
	return false
}
func NewProfile(filename string) (*Profile, error) {
	profile := &Profile{filename: filename}
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, profile)

		if err == nil {
			InitColor(&profile.Colors.Gain, defaultGainColor)
			InitColor(&profile.Colors.Loss, defaultLossColor)
			InitColor(&profile.Colors.Tag, defaultTagColor)
			InitColor(&profile.Colors.Header, defaultHeaderColor)
			InitColor(&profile.Colors.Time, defaultTimeColor)
			InitColor(&profile.Colors.Default, defaultColor)

			profile.SetFilter(profile.Filter)
		}
	} else {
		profile.InitDefaultProfile()
		err = nil
	}
	profile.selectedColumn = -1

	if profile.UpDownJump < 1 {
		profile.UpDownJump = 10
	}

	return profile, err
}

func (profile *Profile) InitDefaultProfile() {
	// Set the refresh intervals to every 3 seconds
	profile.MarketRefresh = 3 // Market data gets fetched every 3 seconds.
	profile.QuotesRefresh = 3 // Stock quotes get updated every 3 seconds.
	profile.Grouped = false
	profile.Tickers = []string{`AAPL`, `C`, `GOOG`, `IBM`, `KO`, `ORCL`, `V`}
	profile.SortColumn = 0 
	profile.Ascending = true 
	profile.Filter = ""
	profile.UpDownJump = 10
	profile.Colors.Gain = defaultGainColor
	profile.Colors.Loss = defaultLossColor
	profile.Colors.Tag = defaultTagColor
	profile.Colors.Header = defaultHeaderColor
	profile.Colors.Time = defaultTimeColor
	profile.Colors.Default = defaultColor
	profile.ShowTimestamp = false
	profile.Save()
}
func InitColor(color *string, defaultValue string) {
	*color = strings.ToLower(*color)
	if !IsSupportedColor(*color) {
		*color = defaultValue
	}
}
func (profile *Profile) Save() error {
	data, err := json.MarshalIndent(profile, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(profile.filename, data, 0644)
}
func (profile *Profile) AddTickers(tickers []string) (added int, err error) {
	added, err = 0, nil
	existing := make(map[string]bool)
	for _, ticker := range profile.Tickers {
		existing[ticker] = true
	}
	for _, ticker := range tickers {
		if _, found := existing[ticker]; !found {
			profile.Tickers = append(profile.Tickers, ticker)
			added++
		}
	}

	if added > 0 {
		sort.Strings(profile.Tickers)
		err = profile.Save()
	}

	return
}
func (profile *Profile) RemoveTickers(tickers []string) (removed int, err error) {
	removed, err = 0, nil
	for _, ticker := range tickers {
		for i, existing := range profile.Tickers {
			if ticker == existing {
				profile.Tickers = append(profile.Tickers[:i], profile.Tickers[i+1:]...)
				removed++
			}
		}
	}

	if removed > 0 {
		err = profile.Save()
	}

	return
}
func (profile *Profile) Reorder() error {
	if profile.selectedColumn == profile.SortColumn {
		profile.Ascending = !profile.Ascending 
	} else {
		profile.SortColumn = profile.selectedColumn 
	}
	return profile.Save()
}
func (profile *Profile) Regroup() error {
	profile.Grouped = !profile.Grouped
	return profile.Save()
}
func (profile *Profile) SetFilter(filter string) {
	if len(filter) > 0 {
		var err error
		profile.filterExpression, err = govaluate.NewEvaluableExpression(filter)

		if err != nil {
			panic(err)
		}

	} else if len(filter) == 0 && profile.filterExpression != nil {
		profile.filterExpression = nil
	}

	profile.Filter = filter
}

func (profile *Profile) ToggleTimestamp() error {
	profile.ShowTimestamp = !profile.ShowTimestamp
	return profile.Save()
}
