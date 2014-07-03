package forms

import (
	"fmt"
	"testing"

	"./elements"
)

func formTest(t *testing.T) {
	rf := NewForm()
	username := elements.NewText()
	username.SetAttribute("name", "username")
	username.SetAttribute("id", "username")
	rf.AddElement(username)

	button := elements.NewSubmit()
	rf.AddElement(rf)

	fmt.Println(rf.Render())
}
