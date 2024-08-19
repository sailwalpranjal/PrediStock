package mop

import (
	"regexp"
	"strings"

	"github.com/nsf/termbox-go"
)

type Markup struct {
	Foreground   termbox.Attribute
	Background   termbox.Attribute
	RightAligned bool
	tags         map[string]termbox.Attribute
	regex        *regexp.Regexp
}

func NewMarkup(profile *Profile) *Markup {
	markup := &Markup{}

	markup.tags = make(map[string]termbox.Attribute)
	markup.tags[`/`] = termbox.ColorDefault
	markup.tags[`black`] = termbox.ColorBlack
	markup.tags[`red`] = termbox.ColorRed
	markup.tags[`green`] = termbox.ColorGreen
	markup.tags[`yellow`] = termbox.ColorYellow
	markup.tags[`blue`] = termbox.ColorBlue
	markup.tags[`magenta`] = termbox.ColorMagenta
	markup.tags[`cyan`] = termbox.ColorCyan
	markup.tags[`white`] = termbox.ColorWhite
	markup.tags[`darkgray`] = termbox.ColorDarkGray
	markup.tags[`lightred`] = termbox.ColorLightRed
	markup.tags[`lightgreen`] = termbox.ColorLightGreen
	markup.tags[`lightyellow`] = termbox.ColorLightYellow
	markup.tags[`lightblue`] = termbox.ColorLightBlue
	markup.tags[`lightmagenta`] = termbox.ColorLightMagenta
	markup.tags[`lightcyan`] = termbox.ColorLightCyan
	markup.tags[`lightgray`] = termbox.ColorLightGray

	markup.tags[`right`] = termbox.ColorDefault
	markup.tags[`b`] = termbox.AttrBold
	markup.tags[`u`] = termbox.AttrUnderline
	markup.tags[`r`] = termbox.AttrReverse

	markup.tags[`gain`] = markup.tags[profile.Colors.Gain]
	markup.tags[`loss`] = markup.tags[profile.Colors.Loss]
	markup.tags[`tag`] = markup.tags[profile.Colors.Tag]
	markup.tags[`header`] = markup.tags[profile.Colors.Header]
	markup.tags[`time`] = markup.tags[profile.Colors.Time]
	markup.tags[`default`] = markup.tags[profile.Colors.Default]

	markup.Foreground = markup.tags[profile.Colors.Default]

	markup.Background = termbox.ColorDefault
	markup.RightAligned = false

	markup.regex = markup.supportedTags()

	return markup
}
func (markup *Markup) Tokenize(str string) []string {
	matches := markup.regex.FindAllStringIndex(str, -1)
	strings := make([]string, 0, len(matches))

	head, tail := 0, 0
	for _, match := range matches {
		tail = match[0]
		if match[1] != 0 {
			if head != 0 || tail != 0 {
				strings = append(strings, str[head:tail])
			}
			strings = append(strings, str[match[0]:match[1]])
		}
		head = match[1]
	}

	if head != len(str) && tail != len(str) {
		strings = append(strings, str[head:])
	}

	return strings
}
func (markup *Markup) IsTag(str string) bool {
	tag, open := probeForTag(str)

	if tag == `` {
		return false
	}

	return markup.process(tag, open)
}

// -----------------------------------------------------------------------------
func (markup *Markup) process(tag string, open bool) bool {
	if attribute, ok := markup.tags[tag]; ok {
		switch tag {
		case `right`:
			markup.RightAligned = open
		default:
			if open {
				if attribute >= termbox.AttrBold {
					markup.Foreground |= attribute
				} else {
					markup.Foreground = attribute
				}
			} else {
				if attribute >= termbox.AttrBold {
					markup.Foreground &= ^attribute
				} else {
					markup.Foreground = markup.tags[`default`]
				}
			}
		}
	}

	return true
}
func (markup *Markup) supportedTags() *regexp.Regexp {
	arr := []string{}

	for tag := range markup.tags {
		arr = append(arr, `</?`+tag+`>`)
	}

	return regexp.MustCompile(strings.Join(arr, `|`))
}

// -----------------------------------------------------------------------------
func probeForTag(str string) (string, bool) {
	if len(str) > 2 && str[0:1] == `<` && str[len(str)-1:] == `>` {
		return extractTagName(str), str[1:2] != `/`
	}

	return ``, false
}
func extractTagName(str string) string {
	if len(str) < 3 {
		return ``
	} else if str[1:2] != `/` {
		return str[1 : len(str)-1]
	} else if len(str) > 3 {
		return str[2 : len(str)-1]
	}

	return `/`
}
