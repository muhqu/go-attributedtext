package ansistyle

import (
	"fmt"
	"strings"

	"github.com/muhqu/go-attributedtext"
	"github.com/muhqu/go-attributedtext/textstyle"
)

func Formater(attr ...attributedtext.Attribute) func(string) string {
	style := AnsiStyleFromAttributes(attr...)
	return func(s string) string {
		return style.Format(s)
	}
}

func Format(as *attributedtext.AttributedString) string {
	sl := as.SliceOfStringsWithAtributes()
	str := ""
	for _, s := range sl {
		style := AnsiStyleFromAttributes(s.Attrs...)
		str += style.Format(s.Str)
	}
	return str
}

type AnsiStyle struct {
	ForegroundColorCode int
	BackgroundColorCode int
	Bold                bool
	Italic              bool
	Underline           bool
}

func AnsiStyleFromAttributes(attrs ...attributedtext.Attribute) *AnsiStyle {
	a := &AnsiStyle{}
	for _, attr := range attrs {
		switch attr := attr.(type) {
		case *textstyle.ForegroundColor:
			if c := ansiColorCode(attr.NamedColor); c >= 0 {
				a.ForegroundColorCode = c + 30
			}
		case *textstyle.BackgroundColor:
			if c := ansiColorCode(attr.NamedColor); c >= 0 {
				a.BackgroundColorCode = c + 40
			}
		case textstyle.TextStyle:
			if attr.IsBold() {
				a.Bold = true
			}
			if attr.IsItalic() {
				a.Italic = true
			}
			if attr.IsUnderline() {
				a.Underline = true
			}
		}
	}
	return a
}

func ansiColorCode(color textstyle.NamedColor) int {
	switch color {
	case textstyle.Black:
		return 0
	case textstyle.Red:
		return 1
	case textstyle.Green:
		return 2
	case textstyle.Yellow:
		return 3
	case textstyle.Blue:
		return 4
	case textstyle.Magenta:
		return 5
	case textstyle.Cyan:
		return 6
	case textstyle.White:
		return 7
	}
	return -1
}

func (a *AnsiStyle) Format(s string) string {
	seq := []string{}
	if a.Bold {
		seq = append(seq, "1")
	}
	if a.Italic {
		seq = append(seq, "3")
	}
	if a.Underline {
		seq = append(seq, "4")
	}
	if a.ForegroundColorCode > 0 {
		seq = append(seq, fmt.Sprintf("%d", a.ForegroundColorCode))
	}
	if a.BackgroundColorCode > 0 {
		seq = append(seq, fmt.Sprintf("%d", a.BackgroundColorCode))
	}
	if len(seq) > 0 {
		s = fmt.Sprintf("\033[%sm%s\033[0m", strings.Join(seq, ";"), s)
	}
	return s
}
