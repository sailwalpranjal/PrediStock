package mop

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

type Screen struct {
	width      int
	height     int
	cleared    bool
	layout     *Layout
	markup     *Markup
	pausedAt   *time.Time
	offset     int
	headerLine int
	max        int
}

func NewScreen(profile *Profile) *Screen {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	screen := &Screen{}
	screen.layout = NewLayout()
	screen.markup = NewMarkup(profile)
	screen.offset = 0

	return screen.Resize()
}

func (screen *Screen) Close() *Screen {
	termbox.Close()

	return screen
}

func (screen *Screen) Resize() *Screen {
	screen.width, screen.height = termbox.Size()
	screen.cleared = false

	return screen
}

func (screen *Screen) Pause(pause bool) *Screen {
	if pause {
		screen.pausedAt = new(time.Time)
		*screen.pausedAt = time.Now()
	} else {
		screen.pausedAt = nil
	}

	return screen
}

func (screen *Screen) Clear() *Screen {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	screen.cleared = true

	return screen
}

func (screen *Screen) ClearLine(x int, y int) *Screen {
	for i := x; i < screen.width; i++ {
		termbox.SetCell(i, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()

	return screen
}

func (screen *Screen) IncreaseOffset(n int) {
	if screen.offset+n <= screen.max {
		screen.offset += n
	} else if screen.max > screen.height {
		screen.offset = screen.max
	}
}

func (screen *Screen) DecreaseOffset(n int) {
	if screen.offset > n {
		screen.offset -= n
	} else {
		screen.offset = 0
	}
}

func (screen *Screen) ScrollTop() {
	screen.offset = 0
}

func (screen *Screen) ScrollBottom() {
	if screen.max > screen.height {
		screen.offset = screen.max
	}
}

func (screen *Screen) DrawOldQuotes(quotes *Quotes) {
	screen.draw(screen.layout.Quotes(quotes), true)
	termbox.Flush()
}

func (screen *Screen) DrawOldMarket(market *Market) {
	screen.draw(screen.layout.Market(market), false)
	termbox.Flush()
}

func (screen *Screen) Draw(objects ...interface{}) *Screen {
	zonename, _ := time.Now().In(time.Local).Zone()
	if screen.pausedAt != nil {
		defer screen.DrawLine(0, 0, `<right><r>`+screen.pausedAt.Format(`3:04:05pm `+zonename)+`</r></right>`)
	}
	for _, ptr := range objects {
		switch ptr.(type) {
		case *Market:
			object := ptr.(*Market)
			screen.draw(screen.layout.Market(object.Fetch()), false)
		case *Quotes:
			object := ptr.(*Quotes)
			screen.draw(screen.layout.Quotes(object.Fetch()), true)
		case time.Time:
			timestamp := ptr.(time.Time).Format(`3:04:05pm ` + zonename)
			screen.DrawLineInverted(0, 0, `<right><time>`+timestamp+`</></right>`)
		default:
			screen.draw(ptr.(string), false)
		}
	}

	termbox.Flush()

	return screen
}
func (screen *Screen) DrawLine(x int, y int, str string) {
	screen.DrawLineFlush(x, y, str, true)
}

func (screen *Screen) DrawLineInverted(x int, y int, str string) {
	screen.DrawLineFlushInverted(x, y, str, true)
}

func (screen *Screen) DrawLineFlush(x int, y int, str string, flush bool) {
	start, column := 0, 0

	for _, token := range screen.markup.Tokenize(str) {
		if screen.markup.IsTag(token) {
			continue
		}
		for i, char := range token {
			if !screen.markup.RightAligned {
				start = x + column
				column++
			} else {
				start = screen.width - len(token) + i
			}
			termbox.SetCell(start, y, char, screen.markup.Foreground, screen.markup.Background)
		}
	}
	if flush {
		termbox.Flush()
	}
}

func (screen *Screen) DrawLineFlushInverted(x int, y int, str string, flush bool) {
	start, column := 0, 0

	for _, token := range screen.markup.Tokenize(str) {
		if screen.markup.IsTag(token) {
			continue
		}
		for i, char := range token {
			if !screen.markup.RightAligned {
				start = x + column
				column++
			} else {
				start = screen.width - len(token) + i
			}
			termbox.SetCell(start, y, char, screen.markup.tags[`black`], screen.markup.Foreground)
		}
	}
	if flush {
		termbox.Flush()
	}
}
func (screen *Screen) draw(str string, offset bool) {
	if !screen.cleared {
		screen.Clear()
	}
	var allLines []string
	drewHeading := false

	screen.width, screen.height = termbox.Size()

	tempFormat := "%" + strconv.Itoa(screen.width) + "s"
	blankLine := fmt.Sprintf(tempFormat, "")
	allLines = strings.Split(str, "\n")

	if offset {
		screen.max = len(allLines) - screen.height + screen.headerLine
	}

	for row := 0; row < len(allLines); row++ {
		if offset {
			if !drewHeading {
				if strings.Contains(allLines[row], "Ticker") &&
					strings.Contains(allLines[row], "Last") &&
					strings.Contains(allLines[row], "Change") {
					drewHeading = true
					screen.headerLine = row
					screen.DrawLine(0, row, allLines[row])
					row += screen.offset
				}
			} else {
				if row <= len(allLines) &&
					row > screen.headerLine {
					screen.DrawLineFlush(0, row-screen.offset, allLines[row], false)
				} else if row > len(allLines)+1 {
					row = len(allLines)
				}
			}
		} else {
			screen.DrawLineFlush(0, row, allLines[row], false)
		}
	}
	if drewHeading {
		for i := len(allLines) - 1 - screen.offset; i < screen.height; i++ {
			if i > screen.headerLine {
				screen.DrawLine(0, i, blankLine)
			}
		}
	}
}
