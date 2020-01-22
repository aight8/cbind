package cbind

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell"
)

type eventHandler func(ev *tcell.EventKey) *tcell.EventKey

// Configuration processes key events by mapping keys to event handlers.
type Configuration struct {
	handlers map[string]eventHandler
	mutex    *sync.RWMutex
}

// NewConfiguration returns a new input configuration.
func NewConfiguration() *Configuration {
	c := Configuration{
		handlers: make(map[string]eventHandler),
		mutex:    new(sync.RWMutex),
	}

	return &c
}

// SetKey sets the handler for a key.
func (c *Configuration) SetKey(mod tcell.ModMask, key tcell.Key, handler func(ev *tcell.EventKey) *tcell.EventKey) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlers[fmt.Sprintf("%d-%d", mod, key)] = handler
}

// SetRune sets the handler for a rune.
func (c *Configuration) SetRune(mod tcell.ModMask, ch rune, handler func(ev *tcell.EventKey) *tcell.EventKey) {
	// Some runes are identical to named keys. Set the bind on the matching
	// named key instead.
	switch ch {
	case '\t':
		c.SetKey(mod, tcell.KeyTab, handler)
		return
	case '\n':
		c.SetKey(mod, tcell.KeyEnter, handler)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlers[fmt.Sprintf("%d:%d", mod, ch)] = handler
}

// Capture handles key events.
func (c *Configuration) Capture(ev *tcell.EventKey) *tcell.EventKey {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if ev == nil {
		return nil
	}

	var keyName string
	if ev.Key() != tcell.KeyRune {
		keyName = fmt.Sprintf("%d-%d", ev.Modifiers(), ev.Key())
	} else {
		keyName = fmt.Sprintf("%d:%d", ev.Modifiers(), ev.Rune())
	}

	handler := c.handlers[keyName]
	if handler != nil {
		return handler(ev)
	}

	return ev
}
