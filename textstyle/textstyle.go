package textstyle

import (
	"strings"

	"github.com/muhqu/go-attributedtext"
)

type TextStyle int

const (
	Bold TextStyle = 1 << iota
	Italic
	Underline
)

func (t TextStyle) IsBold() bool {
	return t&Bold == Bold
}
func (t TextStyle) IsItalic() bool {
	return t&Italic == Italic
}
func (t TextStyle) IsUnderline() bool {
	return t&Underline == Underline
}

func (t TextStyle) Flatten() attributedtext.Attributes {
	styles := make(attributedtext.Attributes, 0)
	if t&Bold == Bold {
		styles = append(styles, Bold)
	}
	if t&Italic == Italic {
		styles = append(styles, Italic)
	}
	if t&Underline == Underline {
		styles = append(styles, Underline)
	}
	return styles
}

func (t TextStyle) String() string {
	styles := make([]string, 0)
	if t&Bold == Bold {
		styles = append(styles, "Bold")
	}
	if t&Italic == Italic {
		styles = append(styles, "Italic")
	}
	if t&Underline == Underline {
		styles = append(styles, "Underline")
	}
	return "TextStyle(" + strings.Join(styles, ",") + ")"
}

type NamedColor int

const (
	Black NamedColor = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	NoColor = -1
)

func (n NamedColor) ColorName() string {
	switch n {
	case Black:
		return "Black"
	case Red:
		return "Red"
	case Green:
		return "Green"
	case Yellow:
		return "Yellow"
	case Blue:
		return "Blue"
	case Magenta:
		return "Magenta"
	case Cyan:
		return "Cyan"
	case White:
		return "White"
	}
	return "NoColor"
}

type ForegroundColor struct {
	NamedColor NamedColor
}

func (f *ForegroundColor) String() string {
	return "ForegroundColor(" + f.NamedColor.ColorName() + ")"
}

type BackgroundColor struct {
	NamedColor NamedColor
}

func (f *BackgroundColor) String() string {
	return "BackgroundColor(" + f.NamedColor.ColorName() + ")"
}

// --------------------------------------------------------

func (a NamedColor) ForegroundColor() *ForegroundColor {
	return &ForegroundColor{NamedColor: a}
}
func (a NamedColor) BackgroundColor() *BackgroundColor {
	return &BackgroundColor{NamedColor: a}
}
