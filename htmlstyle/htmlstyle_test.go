package htmlstyle_test

import (
	"fmt"
	"strings"

	"go/scanner"
	"go/token"

	"github.com/muhqu/go-attributedtext"
	"github.com/muhqu/go-attributedtext/textstyle"

	"github.com/muhqu/go-attributedtext/htmlstyle"
)

func ExampleFormater() {

	err := htmlstyle.Formater(textstyle.Red.ForegroundColor())
	info := htmlstyle.Formater(textstyle.Yellow.ForegroundColor())
	fatal := htmlstyle.Formater(textstyle.Red.ForegroundColor(), textstyle.Bold, textstyle.Underline, textstyle.White.BackgroundColor())

	fmt.Printf("%s\n", err("Hello"))
	fmt.Printf("%s\n", info("World"))
	fmt.Printf("%s\n", fatal("<Fatal>"))

	// Output:
	//
	// <font color="red">Hello</font>
	// <font color="yellow">World</font>
	// <font bgcolor="white"><u><b><font color="red">&lt;Fatal&gt;</font></b></u></font>
}

// Highlight golang source code
func Example_2() {
	src := `
func Hello (s string) string { // go Hello World!
	return fmt.Sprintf("Hello %s", s)
}
`

	type someToken int
	const (
		FUNC_IDENT someToken = iota
		PERIOD_IDENT
	)
	goattributor := attributedtext.NewAttributor(func(as *attributedtext.AttributedString) *attributedtext.AttributedString {
		bytes := ([]byte)(as.Str)
		var s scanner.Scanner
		fset := token.NewFileSet()                        // positions are relative to fset
		file := fset.AddFile("", fset.Base(), len(bytes)) // register input "file"
		s.Init(file, bytes, nil /* no error handler */, scanner.ScanComments)
		lasttok := (interface{})(nil)
		for {
			pos, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			l := len(lit)
			if l == 0 {
				l = len(tok.String())
			}
			if tok == token.IDENT && lasttok == token.FUNC {
				as.AddAttribute(fset.Position(pos).Offset, l, FUNC_IDENT)
			} else if tok == token.IDENT && lasttok == token.PERIOD {
				as.AddAttribute(fset.Position(pos).Offset, l, PERIOD_IDENT)
			} else {
				as.AddAttribute(fset.Position(pos).Offset, l, tok)
			}
			lasttok = tok
		}
		return as
	})
	mapping := attributedtext.AttributeMapping{
		token.FUNC:    attributedtext.Multi(textstyle.Yellow.ForegroundColor(), textstyle.Bold),
		FUNC_IDENT:    attributedtext.Multi(textstyle.Yellow.ForegroundColor()),
		PERIOD_IDENT:  attributedtext.Multi(textstyle.Magenta.ForegroundColor()),
		token.RETURN:  attributedtext.Multi(textstyle.Yellow.ForegroundColor(), textstyle.Bold),
		token.IDENT:   attributedtext.Multi(textstyle.White.ForegroundColor(), textstyle.Bold),
		token.INT:     attributedtext.Multi(textstyle.Green.ForegroundColor()),
		token.LBRACE:  attributedtext.Multi(textstyle.White.ForegroundColor(), textstyle.Bold),
		token.RBRACE:  attributedtext.Multi(textstyle.White.ForegroundColor(), textstyle.Bold),
		token.LPAREN:  attributedtext.Multi(textstyle.White.ForegroundColor(), textstyle.Bold),
		token.RPAREN:  attributedtext.Multi(textstyle.White.ForegroundColor(), textstyle.Bold),
		token.COMMENT: attributedtext.Multi(textstyle.Black.ForegroundColor(), textstyle.Bold),
		token.STRING:  attributedtext.Multi(textstyle.Green.ForegroundColor()),
	}
	mappingattributor := attributedtext.NewAttributeMappingAttributor(mapping)
	attributor := attributedtext.NewMultiAttributor(goattributor, mappingattributor)

	attributedStr := attributor.AttributeString(src)
	htmlStr := htmlstyle.Format(attributedStr)

	for _, line := range strings.Split(htmlStr, "\n") {
		fmt.Printf("%s\n", line)
	}

	// Output:
	//
	// <br>
	// <b><font color="yellow">func</font></b> <font color="yellow">Hello</font> <b><font color="white">(</font></b><b><font color="white">s</font></b> <b><font color="white">string</font></b><b><font color="white">)</font></b> <b><font color="white">string</font></b> <b><font color="white">{</font></b> <b><font color="black">// go Hello World!</font></b><br>
	// &nbsp;&nbsp;&nbsp;&nbsp;<b><font color="yellow">return</font></b> <b><font color="white">fmt</font></b>.<font color="magenta">Sprintf</font><b><font color="white">(</font></b><font color="green">&#34;Hello %s&#34;</font>, <b><font color="white">s</font></b><b><font color="white">)</font></b><br>
	// <b><font color="white">}</font></b><br>
	//
}
