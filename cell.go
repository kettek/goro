package goro

// Cell abstractly represents a rune, style, and other data.
type Cell struct {
	Rune         rune
	Style        Style
	Redraw       bool
	Dirty        bool
	PendingRune  rune
	PendingStyle Style
}
