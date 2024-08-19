package mop

import "github.com/nsf/termbox-go"

type ColumnEditor struct {
	screen  *Screen
	quotes  *Quotes
	layout  *Layout
	profile *Profile
}

func NewColumnEditor(screen *Screen, quotes *Quotes) *ColumnEditor {
	editor := &ColumnEditor{
		screen:  screen,
		quotes:  quotes,
		layout:  screen.layout,
		profile: quotes.profile,
	}

	editor.selectCurrentColumn()

	return editor
}
func (editor *ColumnEditor) Handle(event termbox.Event) bool {
	defer editor.redrawHeader()

	switch event.Key {
	case termbox.KeyEsc:
		return editor.done()

	case termbox.KeyEnter:
		editor.execute()

	case termbox.KeyArrowLeft:
		editor.selectLeftColumn()

	case termbox.KeyArrowRight:
		editor.selectRightColumn()
	}

	return false
}
//-----------------------------------------------------------------------------
func (editor *ColumnEditor) selectCurrentColumn() *ColumnEditor {
	editor.profile.selectedColumn = editor.profile.SortColumn
	editor.redrawHeader()
	return editor
}
//-----------------------------------------------------------------------------
func (editor *ColumnEditor) selectLeftColumn() *ColumnEditor {
	editor.profile.selectedColumn--
	if editor.profile.selectedColumn < 0 {
		editor.profile.selectedColumn = editor.layout.TotalColumns() - 1
	}
	return editor
}
//-----------------------------------------------------------------------------
func (editor *ColumnEditor) selectRightColumn() *ColumnEditor {
	editor.profile.selectedColumn++
	if editor.profile.selectedColumn > editor.layout.TotalColumns()-1 {
		editor.profile.selectedColumn = 0
	}
	return editor
}
//-----------------------------------------------------------------------------
func (editor *ColumnEditor) execute() *ColumnEditor {
	if editor.profile.Reorder() == nil {
		editor.screen.Draw(editor.quotes)
	}

	return editor
}
//-----------------------------------------------------------------------------
func (editor *ColumnEditor) done() bool {
	editor.profile.selectedColumn = -1
	return true
}
//-----------------------------------------------------------------------------
func (editor *ColumnEditor) redrawHeader() {
	editor.screen.DrawLine(0, 4, editor.layout.Header(editor.profile))
	termbox.Flush()
}
