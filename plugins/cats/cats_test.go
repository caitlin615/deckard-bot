package cats

import (
	"fmt"
)

func ExamplePlugin_HandleMessage() {
	fmt.Printf("%q\n", reCats.FindStringSubmatch("!cat"))
	fmt.Printf("%q\n", reCatsType.FindStringSubmatch("!cat fact")[1])
	fmt.Printf("%q\n", reCatsType.FindStringSubmatch("!cat gif"))
	fmt.Printf("%q\n", reCats.FindStringSubmatch("cat"))
	// Output:
	// ["!cat"]
	// "fact"
	// ["!cat gif" "gif"]
	// []
}
