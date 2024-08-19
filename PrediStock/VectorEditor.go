package mop

import (
	"regexp"
	"strings"

	"github.com/nsf/termbox-go"
)
type LineEditor struct {
	command rune          
	cursor  int            
	prompt  string        
	input   string        
	screen  *Screen       
	quotes  *Quotes       
	regex   *regexp.Regexp 
}
func NewLineEditor(screen *Screen, quotes *Quotes) *LineEditor {
	return &LineEditor{
		screen: screen,
		quotes: quotes,
		regex:  regexp.MustCompile(`[,\s]+`),
	}
}
func (editor *LineEditor) Prompt(command rune) *LineEditor {
	filterPrompt := `Set filter: `

	if filter := editor.quotes.profile.Filter; len(filter) > 0 {
		filterPrompt = `Set filter (` + filter + `): `
	}

	prompts := map[rune]string{
		'+': `Add tickers: `, '-': `Remove tickers: `,
		'f': filterPrompt,
	}
	if prompt, ok := prompts[command]; ok {
		editor.prompt = prompt
		editor.command = command

		editor.screen.DrawLine(0, 3, `<white>`+editor.prompt+`</>`)
		termbox.SetCursor(len(editor.prompt), 3)
		termbox.Flush()
	}

	return editor
}
func (editor *LineEditor) Handle(ev termbox.Event) bool {
	defer termbox.Flush()

	switch ev.Key {
	case termbox.KeyEsc:
		return editor.done()

	case termbox.KeyEnter:
		return editor.execute().done()

	case termbox.KeyBackspace, termbox.KeyBackspace2:
		editor.deletePreviousCharacter()

	case termbox.KeyCtrlB, termbox.KeyArrowLeft:
		editor.moveLeft()

	case termbox.KeyCtrlF, termbox.KeyArrowRight:
		editor.moveRight()

	case termbox.KeyCtrlA:
		editor.jumpToBeginning()

	case termbox.KeyCtrlE:
		editor.jumpToEnd()

	case termbox.KeySpace:
		editor.insertCharacter(' ')

	default:
		if ev.Ch != 0 {
			editor.insertCharacter(ev.Ch)
		}
	}

	return false
}

// -----------------------------------------------------------------------------
func (editor *LineEditor) deletePreviousCharacter() *LineEditor {
	if editor.cursor > 0 {
		if editor.cursor < len(editor.input) {
			editor.input = editor.input[0:editor.cursor-1] + editor.input[editor.cursor:len(editor.input)]
		} else {
			editor.input = editor.input[:len(editor.input)-1]
		}
		editor.screen.DrawLine(len(editor.prompt), 3, editor.input+` `) 
		editor.moveLeft()
	}

	return editor
}

// -----------------------------------------------------------------------------
func (editor *LineEditor) insertCharacter(ch rune) *LineEditor {
	if editor.cursor < len(editor.input) {
		editor.input = editor.input[0:editor.cursor] + string(ch) + editor.input[editor.cursor:len(editor.input)]
	} else {
		editor.input += string(ch)
	}
	editor.screen.DrawLine(len(editor.prompt), 3, editor.input)
	editor.moveRight()

	return editor
}

// -----------------------------------------------------------------------------
func (editor *LineEditor) moveLeft() *LineEditor {
	if editor.cursor > 0 {
		editor.cursor--
		termbox.SetCursor(len(editor.prompt)+editor.cursor, 3)
	}

	return editor
}

// -----------------------------------------------------------------------------
func (editor *LineEditor) moveRight() *LineEditor {
	if editor.cursor < len(editor.input) {
		editor.cursor++
		termbox.SetCursor(len(editor.prompt)+editor.cursor, 3)
	}

	return editor
}

// -----------------------------------------------------------------------------
func (editor *LineEditor) jumpToBeginning() *LineEditor {
	editor.cursor = 0
	termbox.SetCursor(len(editor.prompt)+editor.cursor, 3)

	return editor
}

// -----------------------------------------------------------------------------
func (editor *LineEditor) jumpToEnd() *LineEditor {
	editor.cursor = len(editor.input)
	termbox.SetCursor(len(editor.prompt)+editor.cursor, 3)

	return editor
}

// -----------------------------------------------------------------------------
func (editor *LineEditor) execute() *LineEditor {
	switch editor.command {
	case '+':
		tickers := editor.tokenize()
		if len(tickers) > 0 {
			if added, _ := editor.quotes.AddTickers(tickers); added > 0 {
				editor.screen.Draw(editor.quotes)
			}
		}
	case '-':
		tickers := editor.tokenize()
		if len(tickers) > 0 {
			before := len(editor.quotes.profile.Tickers)
			if removed, _ := editor.quotes.RemoveTickers(tickers); removed > 0 {
				editor.screen.Draw(editor.quotes)
				after := before - removed
				for i := before + 1; i > after; i-- {
					editor.screen.ClearLine(0, i+4)
				}
			}
		}
	case 'f':
		if len(editor.input) == 0 {
			editor.input = editor.quotes.profile.Filter
		}

		editor.quotes.profile.SetFilter(editor.input)
	case 'F':
		editor.quotes.profile.SetFilter("")
	}

	return editor
}

// -----------------------------------------------------------------------------
func (editor *LineEditor) done() bool {
	editor.screen.ClearLine(0, 3)
	termbox.HideCursor()

	return true
}
func (editor *LineEditor) tokenize() []string {
	input := strings.ToUpper(strings.Trim(editor.input, `, `))
	return editor.regex.Split(input, -1)
}
