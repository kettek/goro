package glyphs

// Type represents the underlying data type of a Glyphs.
type Type uint8

// These are our built-in GlyphTypes.
const (
	BitmapType Type = iota
	TruetypeType
)

// Glyphs is the interface to a font face.
type Glyphs interface {
	Type() Type
	SetSize(float64)
	Width() int
	Height() int
}
