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
// This function creates a `Profile` by reading and parsing the JSON file at `filename`. If the file is valid, it unmarshals the data, initializes color settings, and applies the filter. If reading the file fails, it initializes the profile with default settings. It also ensures `UpDownJump` is set to at least 10. The function returns the profile and any error encountered.
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
// This function initializes the `Profile` with default values: it sets refresh intervals for market data and quotes to 3 seconds, disables grouping, and sets a default list of tickers. It also configures sorting, filtering, and jumping behavior, assigns default colors to various profile attributes, and disables timestamp display. Finally, it saves the initialized profile.
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
// This function takes a pointer to a color string and a default color value. It converts the color string to lowercase and checks if it is a supported color. If the color is not supported, it assigns the default color value to the provided color string.
func InitColor(color *string, defaultValue string) {
	*color = strings.ToLower(*color)
	if !IsSupportedColor(*color) {
		*color = defaultValue
	}
}
// This function serializes the `Profile` object into a formatted JSON string and writes it to the file specified by `profile.filename`. If the serialization fails, it returns an error. Otherwise, it writes the data to the file with appropriate file permissions (`0644`).
func (profile *Profile) Save() error {
	data, err := json.MarshalIndent(profile, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(profile.filename, data, 0644)
}
// This function adds new tickers to the `Profile`'s `Tickers` list, ensuring no duplicates are added. It first creates a map of existing tickers for quick lookup, then appends each unique ticker from the input list. If any tickers are added, the list is sorted, and the profile is saved. The function returns the number of added tickers and any error encountered during the save process.
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
// This function removes specified tickers from the `Profile`'s `Tickers` list. It iterates through the input tickers and removes matching ones from the profile's list. If any tickers are removed, the profile is saved. The function returns the number of removed tickers and any error encountered during the save process.
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
// This function adjusts the sorting order of the `Profile` based on the selected column. If the selected column is the same as the current sort column, it toggles the sorting direction (ascending or descending). Otherwise, it updates the sort column to the selected column. Afterward, it saves the updated profile.
func (profile *Profile) Reorder() error {
	if profile.selectedColumn == profile.SortColumn {
		profile.Ascending = !profile.Ascending 
	} else {
		profile.SortColumn = profile.selectedColumn 
	}
	return profile.Save()
}
// This function toggles the `Grouped` state of the `Profile`, changing it from grouped to ungrouped or vice versa. After updating the state, it saves the profile.
func (profile *Profile) Regroup() error {
	profile.Grouped = !profile.Grouped
	return profile.Save()
}
// This function sets a filter expression for the `Profile`. If a non-empty filter string is provided, it compiles the expression using the `govaluate` package. If the filter is empty and there is an existing filter expression, it clears it. The filter string is then stored in the `Profile`.
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
// This function toggles the `ShowTimestamp` state of the `Profile`, enabling or disabling the display of timestamps. After updating the state, it saves the profile.
func (profile *Profile) ToggleTimestamp() error {
	profile.ShowTimestamp = !profile.ShowTimestamp
	return profile.Save()
}
