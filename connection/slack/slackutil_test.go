package slack

import (
	"fmt"
)

func ExampleformatSlackMsg() {
	fmt.Println(formatSlackMsg("<http://handwriting.io|handwriting.io>"))
	fmt.Println(formatSlackMsg("the brown dog"))
	fmt.Println(formatSlackMsg(""))
	fmt.Println(formatSlackMsg("<@U2934234|caitlin>"))

	// Output:
	// handwriting.io
	// the brown dog
	//
	// <@U2934234|caitlin>
}
