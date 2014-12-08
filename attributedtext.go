package attributedtext

import "fmt"

func NewAttributedString(Str string) *AttributedString {
	return &AttributedString{Str, make([]*AttributeAtRange, 0)}
}

type SliceOfStringsWithAtributes []*StringWithAttributes
type StringWithAttributes struct {
	Str   string
	Attrs Attributes
}

type Attributes []Attribute

func (a Attributes) String() string {
	str := "Attributes("
	for i, attr := range a {
		if i > 0 {
			str += ", "
		}
		str += fmt.Sprintf("%#v", attr)
	}
	str += ")"
	return str
}

type AttributedString struct {
	Str   string
	Attrs []*AttributeAtRange
}

type AttributeAtRange struct {
	Index     int
	Length    int
	Attribute Attribute
}

type Attribute interface{}
type MultiAttribute interface {
	Flatten() Attributes
}

func Multi(a Attribute, b ...Attribute) Attributes {
	f := make([]Attribute, 0)
	f = append(f, a)
	f = append(f, b...)
	return f
}

func (m Attributes) Flatten() Attributes {
	f := make([]Attribute, 0)
	for _, a := range m {
		if ma, ok := a.(MultiAttribute); ok {
			f = append(f, ma.Flatten()...)
		} else {
			f = append(f, a)
		}
	}
	return f
}

func (a *AttributedString) AddAttribute(index, length int, attribute Attribute) {
	a.Attrs = append(a.Attrs, &AttributeAtRange{index, length, attribute})
}

func (a *AttributedString) String() string {
	return a.Str
}

func (a *AttributedString) SliceOfStringsWithAtributes() SliceOfStringsWithAtributes {
	starts := make(map[int][]Attribute)
	ends := make(map[int][]Attribute)
	actives := make([]Attribute, 0)

	for _, sa := range a.Attrs {
		end := sa.Index + sa.Length
		if ma, ok := sa.Attribute.(MultiAttribute); ok {
			starts[sa.Index] = append(starts[sa.Index], ma.Flatten()...)
			ends[end] = append(ends[end], ma.Flatten()...)
		} else {
			starts[sa.Index] = append(starts[sa.Index], sa.Attribute)
			ends[end] = append(ends[end], sa.Attribute)
		}
	}

	sas := make(SliceOfStringsWithAtributes, 0)

	max := len(a.Str)
	istart := 0
	for i := 0; i < max; i++ {
		if starts[i] != nil || ends[i] != nil {
			if istart != i {
				sas = append(sas, &StringWithAttributes{a.Str[istart:i], actives})
			}
			istart = i
			if ends[i] != nil {
				actives = remove(actives, ends[i]...)
			}
			if starts[i] != nil {
				actives = append(actives, starts[i]...)
			}
		}
	}
	if istart < max {
		sas = append(sas, &StringWithAttributes{a.Str[istart:max], actives})
	}
	return sas
}

func remove(s []Attribute, r ...Attribute) []Attribute {
	var p []Attribute // == nil
	for _, v := range s {
		skip := false
		for _, d := range r {
			if v == d {
				skip = true
				break
			}
		}
		if !skip {
			p = append(p, v)
		}
	}
	return p
}
