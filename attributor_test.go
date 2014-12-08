package attributedtext_test

import (
	"fmt"
	"regexp"

	"github.com/muhqu/go-attributedtext"
)

func ExampleNewRegexAttributor() {

	str := "Hello World,\n\"string\" is a string and 123 is a number."

	exp := regexp.MustCompile(`(?P<STR>"(?:[^"]|\\")*")|(?P<NUM>\d+)`)
	attributor := attributedtext.NewRegexAttributor(exp)
	attrstr := attributor.AttributeString(str)

	for i, swa := range attrstr.SliceOfStringsWithAtributes() {
		fmt.Printf("%d. %-20s %#v\n", i+1, "String:", swa.Str)
		for k, attr := range swa.Attrs {
			fmt.Printf("   %-20s %T %+v\n", fmt.Sprintf("Attribute(%d):", k+1), attr, attr)
		}
	}

	// Output:
	//
	// 1. String:              "Hello World,\n"
	// 2. String:              "\"string\""
	//    Attribute(1):        string STR
	// 3. String:              " is a string and "
	// 4. String:              "123"
	//    Attribute(1):        string NUM
	// 5. String:              " is a number."
}
