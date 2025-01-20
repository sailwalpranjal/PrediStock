package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/mop-tracker/mop"
	"github.com/nsf/termbox-go"
)

const defaultProfile = `.moprc`

const help = `
<u>Command</u>    <u>Description                                </u>
   +                  Add stocks to list
   -                  Remove stocks from list
   ? h H              Display this help screen
   f                  Set filtering expression
   F                  Unset filtering expression
   g G                Group stocks by advancing/declining issues
   o                  Change column sort order
   p P                Pause market data and stock updates
   t                  Toggle timestamp on/off
   Mouse Scroll       Scroll up/down
   PgUp/PgDn          Scroll up/down
   Up/Down arrows     Scroll up
   j J                Scroll up
   k K                Scroll down
   q esc              Quit mop

Enter comma-delimited list of stock tickers when prompted.

<r> Press any key to continue </r>
`

// The mainLoop method is responsible for initiating the event loop in a terminal-based application, managing user input through keyboard and mouse and periodically updating data. The profile's intervals for updating market data, quotes, and timestamps are specified by the timers. Screen rendering, market and quote data generation, as well as asynchronous keyboard input in sane goroutine are also handled by the function. Flags are employed to manage display actions like offering help or halting updates.either way.
func mainLoop(screen *mop.Screen, profile *mop.Profile) {
	var lineEditor *mop.LineEditor
	var columnEditor *mop.ColumnEditor
	termbox.SetInputMode(termbox.InputMouse)
	keyboardQueue := make(chan termbox.Event, 128)

	timestampQueue := time.NewTicker(1 * time.Second)
	quotesQueue := time.NewTicker(time.Duration(profile.QuotesRefresh) * time.Second)
	marketQueue := time.NewTicker(time.Duration(profile.MarketRefresh) * time.Second)
	showingHelp := false
	paused := false
	showingTimestamp := profile.ShowTimestamp
	upDownJump := profile.UpDownJump
	redrawQuotesFlag := false
	redrawMarketFlag := false

	go func() {
		for {
			keyboardQueue <- termbox.PollEvent()
		}
	}()

	market := mop.NewMarket()
	quotes := mop.NewQuotes(market, profile)
	screen.Draw(market)
	screen.Draw(quotes)

// It listens for keyboard input, changes in screen size and the use of mouse movements; events are then looped through. It performs certain actions, such as opening editors, visibility of timestamps, pausing, or going through quotes and market data. However, it does not provide specific commands for this. The editor's input is handled by the active editor, but it can also update the screen at specific intervals using timers, updating quotes and market data. The loop is responsible for handling terminal resizing and scrolling events with the mouse. By taking into account flags and user interactions, the screen is redrawn when necessary, while still adhering to the progress bars and help state.
loop:
	for {
		select {
		case event := <-keyboardQueue:
			switch event.Type {
			case termbox.EventKey:
				if lineEditor == nil && columnEditor == nil && !showingHelp {
					if event.Key == termbox.KeyEsc || event.Ch == 'q' || event.Ch == 'Q' {
						break loop
					} else if event.Ch == '+' || event.Ch == '-' {
						lineEditor = mop.NewLineEditor(screen, quotes)
						lineEditor.Prompt(event.Ch)
					} else if event.Ch == 'f' {
						lineEditor = mop.NewLineEditor(screen, quotes)
						lineEditor.Prompt(event.Ch)
					} else if event.Ch == 'F' {
						profile.SetFilter("")
					} else if event.Ch == 'o' || event.Ch == 'O' {
						columnEditor = mop.NewColumnEditor(screen, quotes)
					} else if event.Ch == 'g' || event.Ch == 'G' {
						if profile.Regroup() == nil {
							screen.Draw(quotes)
						}
					} else if event.Ch == 'p' || event.Ch == 'P' {
						paused = !paused
						screen.Pause(paused).Draw(time.Now())
					} else if event.Ch == '?' || event.Ch == 'h' || event.Ch == 'H' {
						showingHelp = true
						screen.Clear().Draw(help)
					} else if event.Key == termbox.KeyPgdn ||
						event.Ch == 'J' {
						screen.IncreaseOffset(upDownJump)
						redrawQuotesFlag = true
					} else if event.Key == termbox.KeyPgup ||
						event.Ch == 'K' {
						screen.DecreaseOffset(upDownJump)
						redrawQuotesFlag = true
					} else if event.Key == termbox.KeyArrowUp || event.Ch == 'k' {
						screen.DecreaseOffset(1)
						redrawQuotesFlag = true
					} else if event.Key == termbox.KeyArrowDown || event.Ch == 'j' {
						screen.IncreaseOffset(1)
						redrawQuotesFlag = true
					} else if event.Key == termbox.KeyHome {
						screen.ScrollTop()
						redrawQuotesFlag = true
					} else if event.Key == termbox.KeyEnd {
						screen.ScrollBottom()
						redrawQuotesFlag = true
					} else if event.Ch == 't' || event.Ch == 'T' {
						if profile.ToggleTimestamp() == nil {
							showingTimestamp = !showingTimestamp
							screen.Clear().Draw(market, quotes)
						}
					}
				} else if lineEditor != nil {
					if done := lineEditor.Handle(event); done {
						lineEditor = nil
					}
				} else if columnEditor != nil {
					if done := columnEditor.Handle(event); done {
						columnEditor = nil
					}
				} else if showingHelp {
					showingHelp = false
					screen.Clear().Draw(market, quotes)
				}
			case termbox.EventResize:
				screen.Resize()
				if !showingHelp {
					redrawQuotesFlag = true
					redrawMarketFlag = true
				} else {
					screen.Draw(help)
				}
			case termbox.EventMouse:
				if lineEditor == nil && columnEditor == nil && !showingHelp {
					switch event.Key {
					case termbox.MouseWheelUp:
						screen.DecreaseOffset(5)
						redrawQuotesFlag = true
					case termbox.MouseWheelDown:
						screen.IncreaseOffset(5)
						redrawQuotesFlag = true
					}
				}
			}

		case <-timestampQueue.C:
			if !showingHelp && !paused && showingTimestamp {
				screen.Draw(time.Now())
			}

		case <-quotesQueue.C:
			if !showingHelp && !paused && len(keyboardQueue) == 0 {
				go quotes.Fetch()
				redrawQuotesFlag = true
			}

		case <-marketQueue.C:
			if !showingHelp && !paused {
				screen.Draw(market)
			}
		}

		if redrawQuotesFlag && len(keyboardQueue) == 0 {
			screen.DrawOldQuotes(quotes)
			redrawQuotesFlag = false
		}
		if redrawMarketFlag && len(keyboardQueue) == 0 {
			screen.Draw(market)
			redrawMarketFlag = false
		}
	}
}

// The `main` function initializes the application by loading the user profile from the specified path, defaulting to the home directory if not provided. It first checks if the profile exists and is valid, and if not, prompts the user to overwrite the corrupted profile with a default one. After successfully loading or initializing the profile, it creates a new screen object and enters the main event loop (`mainLoop`) to start the application. Upon exiting the event loop, the profile is saved to ensure any changes are persisted. If an error occurs during any step, the program will handle it by either panicking or prompting the user for input.
func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	profileName := flag.String("profile", path.Join(usr.HomeDir, defaultProfile), "path to profile")
	flag.Parse()

	profile, err := mop.NewProfile(*profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "The profile read from `%s` is corrupted.\n\tError: %s\n\n", *profileName, err)
		for {
			fmt.Fprintln(os.Stderr, "Do you want to overwrite the current profile with the default one? [y/n]")
			rne, _, _ := keyboard.GetSingleKey()
			res := strings.ToLower(string(rne))
			if res != "y" && res != "n" {
				fmt.Fprintf(os.Stderr, "Invalid answer `%s`\n\n", res)
				continue
			}

			if res == "y" {
				profile.InitDefaultProfile()
				break
			} else {
				os.Exit(1)
			}
		}
	}
	screen := mop.NewScreen(profile)
	defer screen.Close()

	mainLoop(screen, profile)
	profile.Save()
}
