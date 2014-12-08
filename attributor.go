package attributedtext

import "regexp"

type AttributorFn func(*AttributedString) *AttributedString

type Attributor interface {
	Attribute(*AttributedString) *AttributedString
	AttributeString(string) *AttributedString
}

func NewAttributor(fn AttributorFn) Attributor {
	return &attributor{fn}
}

type attributor struct {
	AttributorFn
}

func (a *attributor) Attribute(s *AttributedString) *AttributedString {
	return a.AttributorFn(&AttributedString{s.Str, s.Attrs})
}

func (a *attributor) AttributeString(s string) *AttributedString {
	return a.AttributorFn(NewAttributedString(s))
}

func NewRegexAttributor(r ...*regexp.Regexp) Attributor {
	return NewAttributor(func(s *AttributedString) *AttributedString {
		for _, r := range r {
			names := r.SubexpNames()
			matches := r.FindAllSubmatchIndex(([]byte)(s.Str), -1)
			if matches != nil {
				for _, m := range matches {
					for n, name := range names {
						if name != "" { // we're only interested in named-capture groups
							// matches come in pairs [m1_start, m1_end, m2_start, m2_end, ...]
							match := m[(n * 2):(n*2 + 2)]
							if match[0] != -1 { // and only those that have matches
								s.AddAttribute(match[0], match[1]-match[0], name)
							}
						}
					}
				}
			}
		}
		return s
	})
}

func NewMultiAttributor(m ...Attributor) Attributor {
	return NewAttributor(func(s *AttributedString) *AttributedString {
		for _, a := range m {
			s = a.Attribute(s)
		}
		return s
	})
}

type AttributeMapping map[Attribute]Attribute

func NewAttributeMappingAttributor(m map[Attribute]Attribute) Attributor {
	return NewAttributor(func(s *AttributedString) *AttributedString {
		for _, aar := range s.Attrs {
			if a, found := m[aar.Attribute]; found {
				s.AddAttribute(aar.Index, aar.Length, a)
			}
		}
		return s
	})
}
