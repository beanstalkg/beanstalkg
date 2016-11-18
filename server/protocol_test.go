package server

import (
	"testing"
	"fmt"
)

func TestCommand_Parse(t *testing.T) {
	c := Command{}
	c.Parse("use test")
	fmt.Println(c)
	if c.RawCommand != "use test" {
		t.Failed()
	}
	if c.Name != "use" {
		t.Failed()
	}

	c = Command{}
	done, err := c.Parse("put 1 2 3 4")

	if err != nil {
		t.Failed()
	}

	if done {
		t.Failed()
	}

	if c.Name != "put" || c.Params["ttr"] != "1" {
		t.Failed()
	}

	c.Parse("data")

	if c.Params["data"] != "data" {
		t.Failed()
	}
}
