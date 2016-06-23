package message

import "fmt"

func Example() {
	m1 := new(Basic)
	m1.Text = "I am a string"
	m1.ID = 2
	m1.Finished = true

	fmt.Println(m1.Text)
	fmt.Println(m1.ID)
	fmt.Println(m1.Finished)

	// Output:
	// I am a string
	// 2
	// true
}
