package textstyle_test

import (
	"fmt"

	"github.com/muhqu/go-attributedtext"
	"github.com/muhqu/go-attributedtext/textstyle"
)

func ExampleRed() {

	attrstr := attributedtext.NewAttributedString("Some Message")
	attrstr.AddAttribute(0, 4, &textstyle.ForegroundColor{textstyle.Red})
	attrstr.AddAttribute(5, 7, textstyle.Underline|textstyle.Bold)

	for i, swa := range attrstr.SliceOfStringsWithAtributes() {
		fmt.Printf("%d. %-20s %#v\n", i+1, "String:", swa.Str)
		for k, attr := range swa.Attrs {
			fmt.Printf("   %-20s %T %+v\n", fmt.Sprintf("Attribute(%d):", k+1), attr, attr)
		}
	}

	// Output:
	//
	// 1. String:              "Some"
	//    Attribute(1):        *textstyle.ForegroundColor ForegroundColor(Red)
	// 2. String:              " "
	// 3. String:              "Message"
	//    Attribute(1):        textstyle.TextStyle TextStyle(Bold)
	//    Attribute(2):        textstyle.TextStyle TextStyle(Underline)
}
