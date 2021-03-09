package cbind

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

const pressTimes = 7

func TestConfiguration(t *testing.T) {
	t.Parallel()

	wg := make([]*sync.WaitGroup, len(testCases))

	config := NewConfiguration()
	for i, c := range testCases {
		wg[i] = new(sync.WaitGroup)
		wg[i].Add(pressTimes)

		i := i // Capture
		if c.key != tcell.KeyRune {
			config.SetKey(c.mod, c.key, func(ev *tcell.EventKey) *tcell.EventKey {
				wg[i].Done()
				return nil
			})
		} else {
			config.SetRune(c.mod, c.ch, func(ev *tcell.EventKey) *tcell.EventKey {
				wg[i].Done()
				return nil
			})
		}

	}

	done := make(chan struct{})
	timeout := time.After(5 * time.Second)

	go func() {
		for i := range testCases {
			wg[i].Wait()
		}

		done <- struct{}{}
	}()

	go func() {
		for j := 0; j < pressTimes; j++ {
			for i, c := range testCases {
				k := tcell.NewEventKey(c.key, c.ch, c.mod)
				if k.Key() != c.key {
					t.Fatalf("failed to test capturing keybinds: tcell modified EventKey.Key: expected %d, got %d", c.key, k.Key())
				} else if k.Rune() != c.ch {
					t.Fatalf("failed to test capturing keybinds: tcell modified EventKey.Rune: expected %d, got %d", c.ch, k.Rune())
				} else if k.Modifiers() != c.mod {
					t.Fatalf("failed to test capturing keybinds: tcell modified EventKey.Modifiers: expected %d, got %d", c.mod, k.Modifiers())
				}

				ev := config.Capture(tcell.NewEventKey(c.key, c.ch, c.mod))
				if ev != nil {
					t.Fatalf("failed to test capturing keybinds: failed to register case %d event %d %d %d", i, c.mod, c.key, c.ch)
				}
			}
		}
	}()

	select {
	case <-timeout:
		t.Error("timeout")
	case <-done:
		// Wait at least one second to catch problems before exiting.
		<-time.After(1 * time.Second)
	}
}

// Example of creating and using an input configuration.
func ExampleNewConfiguration() {
	// Create a new input configuration to store the key bindings.
	c := NewConfiguration()

	handleSave := func(ev *tcell.EventKey) *tcell.EventKey {
		// Save
		return nil
	}

	handleOpen := func(ev *tcell.EventKey) *tcell.EventKey {
		// Open
		return nil
	}

	handleExit := func(ev *tcell.EventKey) *tcell.EventKey {
		// Exit
		return nil
	}

	// Bind Alt+s.
	if err := c.Set("Alt+s", handleSave); err != nil {
		log.Fatalf("failed to set keybind: %s", err)
	}

	// Bind Alt+o.
	c.SetRune(tcell.ModAlt, 'o', handleOpen)

	// Bind Escape.
	c.SetKey(tcell.ModNone, tcell.KeyEscape, handleExit)

	// Capture input. This will differ based on the framework in use (if any).
	// When using tview or cview, call Application.SetInputCapture before calling
	// Application.Run.
	// app.SetInputCapture(c.Capture)
}

// Example of capturing key events.
func ExampleConfiguration_Capture() {
	// See the end of the NewConfiguration example.
}
