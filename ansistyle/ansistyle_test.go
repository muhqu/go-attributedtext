package ansistyle_test

import (
	"fmt"
	"regexp"
	"strings"

	"go/scanner"
	"go/token"

	"github.com/muhqu/go-attributedtext"
	"github.com/muhqu/go-attributedtext/ansistyle"
	"github.com/muhqu/go-attributedtext/textstyle"
)

func ExampleFormater() {

	err := ansistyle.Formater(textstyle.Red.ForegroundColor())
	info := ansistyle.Formater(textstyle.Yellow.ForegroundColor())
	fatal := ansistyle.Formater(textstyle.Red.ForegroundColor(), textstyle.Bold, textstyle.Underline, textstyle.White.BackgroundColor())

	fmt.Printf("%#v\n", err("Hello"))
	fmt.Printf("%#v\n", info("World"))
	fmt.Printf("%#v\n", fatal("Fatal!"))

	// Output:
	//
	// "\x1b[31mHello\x1b[0m"
	// "\x1b[33mWorld\x1b[0m"
	// "\x1b[1;4;31;47mFatal!\x1b[0m"
}

func ExampleFormatedString() {

	attrstr := attributedtext.NewAttributedString("Some Message")
	attrstr.AddAttribute(0, 4, &textstyle.ForegroundColor{textstyle.Red})
	attrstr.AddAttribute(5, 7, textstyle.Underline|textstyle.Bold)

	fmt.Printf("%#v\n", ansistyle.Format(attrstr))

	// Output:
	//
	// "\x1b[31mSome\x1b[0m \x1b[1;4mMessage\x1b[0m"
}

// Regex based log attributor/colorer
func Example_1() {
	str := `
[Wed Oct 11 14:32:50 2000] [info] [client 127.0.0.1] "GET /apache_pb.gif HTTP/1.0" 200 2326
[Wed Oct 11 14:32:52 2000] [error] [client 127.0.0.1] client denied by server configuration: /export/home/live/ap/htdocs/test
`

	exp1 := regexp.MustCompile(`(\[(?P<date>[^\]]+)\] ((?P<error>\[error\])|(?P<info>\[info\])))`)
	exp2 := regexp.MustCompile(`(?P<ipaddress>[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})`)
	regexattributor := attributedtext.NewRegexAttributor(exp1, exp2)
	mapping := attributedtext.AttributeMapping{
		"date":      attributedtext.Multi(textstyle.Italic),
		"info":      attributedtext.Multi(textstyle.Green.ForegroundColor(), textstyle.Bold),
		"error":     attributedtext.Multi(textstyle.Red.ForegroundColor(), textstyle.Bold),
		"ipaddress": attributedtext.Multi(textstyle.Italic, textstyle.Cyan.ForegroundColor()),
	}
	mappingattributor := attributedtext.NewAttributeMappingAttributor(mapping)
	attributor := attributedtext.NewMultiAttributor(regexattributor, mappingattributor)

	attributedStr := attributor.AttributeString(str)
	ansiColoredStr := ansistyle.Format(attributedStr)

	for _, line := range strings.Split(ansiColoredStr, "\n") {
		fmt.Printf("%#v\n", line)
	}

	// Output:
	// ""
	// "[\x1b[3mWed Oct 11 14:32:50 2000\x1b[0m] \x1b[1;32m[info]\x1b[0m [client \x1b[3;36m127.0.0.1\x1b[0m] \"GET /apache_pb.gif HTTP/1.0\" 200 2326"
	// "[\x1b[3mWed Oct 11 14:32:52 2000\x1b[0m] \x1b[1;31m[error]\x1b[0m [client \x1b[3;36m127.0.0.1\x1b[0m] client denied by server configuration: /export/home/live/ap/htdocs/test"
	// ""
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
	ansiColoredStr := ansistyle.Format(attributedStr)

	for _, line := range strings.Split(ansiColoredStr, "\n") {
		fmt.Printf("%#v\n", line)
	}

	// Output:
	//
	// ""
	// "\x1b[1;33mfunc\x1b[0m \x1b[33mHello\x1b[0m \x1b[1;37m(\x1b[0m\x1b[1;37ms\x1b[0m \x1b[1;37mstring\x1b[0m\x1b[1;37m)\x1b[0m \x1b[1;37mstring\x1b[0m \x1b[1;37m{\x1b[0m \x1b[1;30m// go Hello World!\x1b[0m"
	// "\t\x1b[1;33mreturn\x1b[0m \x1b[1;37mfmt\x1b[0m.\x1b[35mSprintf\x1b[0m\x1b[1;37m(\x1b[0m\x1b[32m\"Hello %s\"\x1b[0m, \x1b[1;37ms\x1b[0m\x1b[1;37m)\x1b[0m"
	// "\x1b[1;37m}\x1b[0m"
	// ""
}
