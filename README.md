# cbind

Key event handling library for tcell

## Features

- Set `KeyEvent` handlers
- Encode and decode `KeyEvent`s as human-readable strings

## Usage

```go
// Create a new input configuration to store the key bindings.
c := NewConfiguration()

// Define save event handler.
handleSave := func(ev *tcell.EventKey) *tcell.EventKey {
    return nil
}

// Define open event handler.
handleOpen := func(ev *tcell.EventKey) *tcell.EventKey {
    return nil
}

// Define exit event handler.
handleExit := func(ev *tcell.EventKey) *tcell.EventKey {
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
app.SetInputCapture(c.Capture)
```

## Documentation

The utility program `whichkeybind` is available to determine and validate key combinations.

```bash
go get github.com/aight8/cbind/whichkeybind
```

## Support

Please share issues and suggestions [here](https://github.com/aight8/cbind/issues).
