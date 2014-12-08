package attributedtext_test

import (
	"fmt"

	"github.com/muhqu/go-attributedtext"
)

func ExampleSliceOfStringsWithAtributes() {
	attrstr := attributedtext.NewAttributedString("Some Message")
	attrstr.AddAttribute(0, 4, "A")
	attrstr.AddAttribute(5, 7, 1234)

	for i, swa := range attrstr.SliceOfStringsWithAtributes() {
		fmt.Printf("%d. %-20s %#v\n", i+1, "String:", swa.Str)
		for k, attr := range swa.Attrs {
			fmt.Printf("   %-20s %T %+v\n", fmt.Sprintf("Attribute(%d):", k+1), attr, attr)
		}
	}

	// Output:
	//
	// 1. String:              "Some"
	//    Attribute(1):        string A
	// 2. String:              " "
	// 3. String:              "Message"
	//    Attribute(1):        int 1234
}
