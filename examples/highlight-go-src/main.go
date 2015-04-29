package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/muhqu/go-attributedtext/textstyle/ansiformater"
	"github.com/muhqu/go-attributedtext/textstyle/htmlformater"
)

func main() {
	var format string
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [ -format FORMAT ] FILE\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&format, "format", "ansi", "output format: ansi, html")
	flag.Parse()
	args := flag.Args()
	//log.Print(args)

	attributor := Highlighter()

	var input io.Reader
	if len(args) == 0 && !isTerminal(os.Stdin) {
		input = os.Stdin
	} else if len(args) == 1 {
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
			return
		}
		defer file.Close()
		input = file
	} else {
		flag.Usage()
		os.Exit(1)
		return
	}

	src, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatal(err)
		return
	}
	attrsrc := attributor.AttributeString(string(src))

	var outstr string
	if format == "html" {
		outstr = wrapHtmlDoc(htmlformater.Format(attrsrc))
	} else {
		outstr = ansiformater.Format(attrsrc)
	}

	os.Stdout.WriteString(outstr)
}

func isTerminal(fd *os.File) bool {
	fi, _ := fd.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func wrapHtmlDoc(html string) string {
	return fmt.Sprintf(`<html>
<head>
<style>
  body {
    background: black;
    color: silver;
    font-family: Monaco, monospaced;
    font-size: 9pt;
  }
 </style>
</head>
<body>
%s
</body>
</html>`, html)
}
