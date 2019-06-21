package goro

// Event is the Event interface.
type Event interface {
}

// EventScreen represents screen changes such as terminal resizing.
type EventScreen struct {
}

// EventKey represents a keypress event.
type EventKey struct {
	Key   Key
	Rune  rune
	State bool
	Shift bool
	Ctrl  bool
	Alt   bool
	Meta  bool
}

func (e EventKey) HasModifiers() bool {
	return e.Shift || e.Ctrl || e.Alt || e.Meta
}

// EventMouse represents a mouse press event.
type EventMouse struct {
	X, Y   int
	Button int
	State  bool
}

// EventQuit represents a quit event.
type EventQuit struct {
}
