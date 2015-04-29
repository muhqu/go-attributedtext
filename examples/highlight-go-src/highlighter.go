package main

import (
	"go/scanner"
	"go/token"

	"github.com/muhqu/go-attributedtext"
	"github.com/muhqu/go-attributedtext/textstyle"
)

type specialtoken int

const (
	stoken_FUNC_IDENT specialtoken = iota
	stoken_PERIOD_IDENT
)

func Highlighter() attributedtext.Attributor {
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
				as.AddAttribute(fset.Position(pos).Offset, l, stoken_FUNC_IDENT)
			} else if tok == token.IDENT && lasttok == token.PERIOD {
				as.AddAttribute(fset.Position(pos).Offset, l, stoken_PERIOD_IDENT)
			} else {
				as.AddAttribute(fset.Position(pos).Offset, l, tok)
			}
			lasttok = tok
		}
		return as
	})
	mapping := attributedtext.AttributeMapping{
		token.FUNC:          {textstyle.Yellow.ForegroundColor(), textstyle.Bold},
		stoken_FUNC_IDENT:   {textstyle.Yellow.ForegroundColor()},
		stoken_PERIOD_IDENT: {textstyle.Magenta.ForegroundColor()},
		token.RETURN:        {textstyle.Yellow.ForegroundColor(), textstyle.Bold},
		token.IDENT:         {textstyle.White.ForegroundColor(), textstyle.Bold},
		token.INT:           {textstyle.Green.ForegroundColor()},
		token.LBRACE:        {textstyle.White.ForegroundColor(), textstyle.Bold},
		token.RBRACE:        {textstyle.White.ForegroundColor(), textstyle.Bold},
		token.LPAREN:        {textstyle.White.ForegroundColor(), textstyle.Bold},
		token.RPAREN:        {textstyle.White.ForegroundColor(), textstyle.Bold},
		token.COMMENT:       {textstyle.Black.ForegroundColor(), textstyle.Bold},
		token.STRING:        {textstyle.Green.ForegroundColor()},
		token.VAR:           {textstyle.Yellow.ForegroundColor(), textstyle.Bold},
		token.IF:            {textstyle.Yellow.ForegroundColor(), textstyle.Bold},
		token.ELSE:          {textstyle.Yellow.ForegroundColor(), textstyle.Bold},
		token.CONST:         {textstyle.Yellow.ForegroundColor(), textstyle.Bold},
	}
	mappingattributor := attributedtext.NewAttributeMappingAttributor(mapping)
	attributor := attributedtext.NewMultiAttributor(goattributor, mappingattributor)

	return attributor
}
