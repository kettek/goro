package glyphs

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// Truetype is our Truetype data.
type Truetype struct {
	truetype      *truetype.Font
	size          float64
	width, height int
	Normal        font.Face
	// How do we manage bold? italics? etc?
	// bold font.Face
	// italics font.Face
}

// LoadTruetype loads a Face from the provided path.
func LoadTruetype(path string) (Glyphs, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadTruetypeFromBytes(bytes)
}

// LoadTruetypeFromBytes loads a Face from the provided bytes of TTF data.
func LoadTruetypeFromBytes(ttf []byte) (Glyphs, error) {
	tt, err := truetype.Parse(ttf)
	if err != nil {
		return nil, err
	}
	face := &Truetype{
		truetype: tt,
	}
	return face, nil
}

// Type returns FaceTruetype.
func (f *Truetype) Type() Type {
	return TruetypeType
}

// SetSize sets the glyph size to the provided size.
func (f *Truetype) SetSize(size float64) {
	if f.size == size {
		return
	}

	f.rebuild()

	f.size = size
}

// Width gets the width of the font.
func (f *Truetype) Width() int {
	return f.width
}

// Height gets the width of the font.
func (f *Truetype) Height() int {
	return f.height
}

// rebuild rebuilds all the font variants.
func (f *Truetype) rebuild() {
	f.Normal = truetype.NewFace(f.truetype, &truetype.Options{
		Size:    f.size,
		DPI:     72,               // FIXME
		Hinting: font.HintingFull, // FIXME
	})

	metrics := f.Normal.Metrics()
	f.height = (metrics.Ascent + metrics.Descent).Round()

	advance, ok := f.Normal.GlyphAdvance('M')
	if !ok {
		f.width = f.height / 2 // This is not good.
	}
	f.width = advance.Round()
}
