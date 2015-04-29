package htmlformater

import (
	"fmt"
	"html"
	"strings"

	"github.com/muhqu/go-attributedtext"
	"github.com/muhqu/go-attributedtext/textstyle"
)

func Formater(attr ...attributedtext.Attribute) func(string) string {
	tag := htmlTagFromAttributes(attr...)
	return func(s string) string {
		return tag.Wrap(s)
	}
}

func Format(as *attributedtext.AttributedString) string {
	sl := as.SliceOfStringsWithAtributes()
	str := ""
	for _, s := range sl {
		tag := htmlTagFromAttributes(s.Attrs...)
		str += tag.Wrap(s.Str)
	}
	return str
}

func htmlTagFromAttributes(attrs ...attributedtext.Attribute) *tag {
	t := &tag{}
	for _, attr := range attrs {
		switch attr := attr.(type) {
		case *textstyle.ForegroundColor:
			t = &tag{tag: "font", attributes: map[string]string{
				"color": strings.ToLower(attr.NamedColor.ColorName()),
			}, wrapped: t}
		case *textstyle.BackgroundColor:
			t = &tag{tag: "font", attributes: map[string]string{
				"bgcolor": strings.ToLower(attr.NamedColor.ColorName()),
			}, wrapped: t}
		case textstyle.TextStyle:
			if attr.IsBold() {
				t = &tag{tag: "b", wrapped: t}
			}
			if attr.IsItalic() {
				t = &tag{tag: "i", wrapped: t}
			}
			if attr.IsUnderline() {
				t = &tag{tag: "u", wrapped: t}
			}
		}
	}
	return t
}

type tag struct {
	tag        string
	attributes map[string]string
	wrapped    *tag
}

func (t *tag) Wrap(s string) string {
	if t.tag == "" {
		s = html.EscapeString(s)
		s = strings.Replace(s, "\n", "<br>\n", -1)
		s = strings.Replace(s, "\t", "    ", -1)
		s = strings.Replace(s, "  ", "&nbsp;&nbsp;", -1)
		return s
	}
	if t.wrapped != nil {
		s = t.wrapped.Wrap(s)
	}
	a := make([]string, len(t.attributes))
	for k, v := range t.attributes {
		a = append(a, fmt.Sprintf("%s=%q", k, v))
	}
	return fmt.Sprintf("<%s%s>%s</%s>", t.tag, strings.Join(a, " "), s, t.tag)
}
